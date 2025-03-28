

resource "aws_vpc" "blockchain_vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "blockchain-vpc"
  }
}

resource "aws_subnet" "public_subnet_1" {
  vpc_id                  = aws_vpc.blockchain_vpc.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "${var.aws_region}a"
  map_public_ip_on_launch = true

  tags = {
    Name = "blockchain-public-subnet-1"
  }
}

resource "aws_subnet" "public_subnet_2" {
  vpc_id                  = aws_vpc.blockchain_vpc.id
  cidr_block              = "10.0.2.0/24"
  availability_zone       = "${var.aws_region}b"
  map_public_ip_on_launch = true

  tags = {
    Name = "blockchain-public-subnet-2"
  }
}

resource "aws_internet_gateway" "blockchain_igw" {
  vpc_id = aws_vpc.blockchain_vpc.id

  tags = {
    Name = "blockchain-igw"
  }
}

resource "aws_route_table" "blockchain_rtb" {
  vpc_id = aws_vpc.blockchain_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.blockchain_igw.id
  }

  tags = {
    Name = "blockchain-rtb"
  }
}

resource "aws_route_table_association" "rtb_assoc_1" {
  subnet_id      = aws_subnet.public_subnet_1.id
  route_table_id = aws_route_table.blockchain_rtb.id
}

resource "aws_route_table_association" "rtb_assoc_2" {
  subnet_id      = aws_subnet.public_subnet_2.id
  route_table_id = aws_route_table.blockchain_rtb.id
}

resource "aws_security_group" "alb_sg" {
  name        = "blockchain-alb-sg"
  description = "Security group for blockchain ALB"
  vpc_id      = aws_vpc.blockchain_vpc.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "blockchain-alb-sg"
  }
}

resource "aws_security_group" "ecs_tasks_sg" {
  name        = "blockchain-ecs-tasks-sg"
  description = "Security group for blockchain ECS tasks"
  vpc_id      = aws_vpc.blockchain_vpc.id

  ingress {
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.alb_sg.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "blockchain-ecs-tasks-sg"
  }
}

resource "aws_lb" "blockchain_alb" {
  name               = "blockchain-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_sg.id]
  subnets            = [aws_subnet.public_subnet_1.id, aws_subnet.public_subnet_2.id]

  tags = {
    Name = "blockchain-alb"
  }
}

resource "aws_lb_target_group" "blockchain_tg" {
  name        = "blockchain-target-group"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = aws_vpc.blockchain_vpc.id
  target_type = "ip"

  health_check {
    path                = "/health"
    port                = "traffic-port"
    healthy_threshold   = 3
    unhealthy_threshold = 3
    timeout             = 5
    interval            = 30
    matcher             = "200"
  }
}

resource "aws_lb_listener" "blockchain_listener" {
  load_balancer_arn = aws_lb.blockchain_alb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.blockchain_tg.arn
  }
}

resource "aws_ecs_cluster" "blockchain_cluster" {
  name = "blockchain-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

resource "aws_iam_role" "ecs_task_execution_role" {
  name = "blockchain-ecs-task-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_cloudwatch_log_group" "blockchain_logs" {
  name              = "/ecs/blockchain-client"
  retention_in_days = 30
}

resource "aws_ecs_task_definition" "blockchain_task" {
  family                   = "blockchain-client"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  container_definitions = jsonencode([
    {
      name      = "blockchain-client"
      image     = "kaytheog/blockchain-client:latest"
      essential = true
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
          protocol      = "tcp"
        }
      ]
      environment = [
        {
          name  = "PORT"
          value = "8080"
        },
        {
          name  = "RPC_URL"
          value = "https://polygon-rpc.com/"
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.blockchain_logs.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])
}

# Create an ECS service
resource "aws_ecs_service" "blockchain_service" {
  name            = "blockchain-service"
  cluster         = aws_ecs_cluster.blockchain_cluster.id
  task_definition = aws_ecs_task_definition.blockchain_task.arn
  desired_count   = 2
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [aws_subnet.public_subnet_1.id, aws_subnet.public_subnet_2.id]
    security_groups  = [aws_security_group.ecs_tasks_sg.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.blockchain_tg.arn
    container_name   = "blockchain-client"
    container_port   = 8080
  }

  depends_on = [aws_lb_listener.blockchain_listener]
}
