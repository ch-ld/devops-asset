#!/bin/bash

# DevOps Asset Management System 一键部署脚本
# 支持生产环境和开发环境部署

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "检查系统依赖..."
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    # 检查Docker Compose
    if ! docker compose version &> /dev/null; then
        log_error "Docker Compose 未安装或版本过低，请安装 Docker Compose v2.0+"
        exit 1
    fi
    
    # 检查Git
    if ! command -v git &> /dev/null; then
        log_error "Git 未安装，请先安装 Git"
        exit 1
    fi
    
    log_success "系统依赖检查通过"
}

# 创建环境配置文件
setup_environment() {
    log_info "设置环境配置..."
    
    if [ ! -f .env ]; then
        if [ -f .env.template ]; then
            cp .env.template .env
            log_success "已创建 .env 配置文件"
            log_warning "请编辑 .env 文件配置数据库密码、JWT密钥等重要参数"
        else
            log_warning ".env.template 文件不存在，将使用默认配置"
        fi
    else
        log_info ".env 文件已存在，跳过创建"
    fi
}

# 创建必要的目录
create_directories() {
    log_info "创建必要的目录..."
    
    directories=(
        "server/log"
        "docker/ssl"
        "docker/mysql/init"
        "docker/redis"
        "backup"
    )
    
    for dir in "${directories[@]}"; do
        if [ ! -d "$dir" ]; then
            mkdir -p "$dir"
            log_info "创建目录: $dir"
        fi
    done
    
    log_success "目录创建完成"
}

# 生成随机密码
generate_random_password() {
    openssl rand -base64 32 | tr -d "=+/" | cut -c1-25
}

# 生成AES密钥
generate_aes_key() {
    openssl rand -base64 32 | tr -d "=+/" | cut -c1-32
}

# 自动配置环境变量
auto_configure_env() {
    if [ "$AUTO_CONFIG" = "true" ]; then
        log_info "自动生成安全配置..."
        
        # 生成随机密码
        MYSQL_ROOT_PWD=$(generate_random_password)
        MYSQL_PWD=$(generate_random_password)
        REDIS_PWD=$(generate_random_password)
        JWT_SECRET=$(generate_random_password)
        AES_KEY=$(generate_aes_key)
        
        # 更新.env文件
        if [ -f .env ]; then
            sed -i.bak \
                -e "s/MYSQL_ROOT_PASSWORD=.*/MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PWD}/" \
                -e "s/MYSQL_PASSWORD=.*/MYSQL_PASSWORD=${MYSQL_PWD}/" \
                -e "s/REDIS_PASSWORD=.*/REDIS_PASSWORD=${REDIS_PWD}/" \
                -e "s/JWT_SECRET=.*/JWT_SECRET=${JWT_SECRET}/" \
                -e "s/AES_KEY=.*/AES_KEY=${AES_KEY}/" \
                .env
            
            log_success "已自动生成安全配置"
            log_warning "请保存以下重要信息："
            echo "MySQL Root密码: ${MYSQL_ROOT_PWD}"
            echo "MySQL用户密码: ${MYSQL_PWD}"
            echo "Redis密码: ${REDIS_PWD}"
        fi
    fi
}

# 部署生产环境
deploy_production() {
    log_info "开始部署生产环境..."
    
    # 构建和启动服务
    docker compose --profile production build
    docker compose --profile production up -d
    
    # 等待数据库启动
    log_info "等待数据库启动..."
    sleep 30
    
    # 执行数据库迁移
    log_info "执行数据库迁移..."
    docker compose exec server ./server --migrate
    
    log_success "生产环境部署完成！"
    echo ""
    echo "访问地址："
    echo "  前端界面: http://localhost"
    echo "  API接口: http://localhost:8080"
    echo "  API文档: http://localhost:8080/swagger/index.html"
    echo "  健康检查: http://localhost:8080/health"
}

# 部署开发环境
deploy_development() {
    log_info "开始部署开发环境..."
    
    # 构建和启动服务
    docker compose --profile development build
    docker compose --profile development up -d
    
    # 等待数据库启动
    log_info "等待数据库启动..."
    sleep 30
    
    log_success "开发环境部署完成！"
    echo ""
    echo "访问地址："
    echo "  前端界面: http://localhost:3000"
    echo "  API接口: http://localhost:8080"
    echo "  API文档: http://localhost:8080/swagger/index.html"
    echo "  健康检查: http://localhost:8080/health"
    echo ""
    echo "开发模式说明："
    echo "  - 后端支持热重载"
    echo "  - 前端支持热重载"
    echo "  - 数据库和Redis数据持久化"
}

# 健康检查
health_check() {
    log_info "执行健康检查..."
    
    # 检查服务状态
    if docker compose ps --services --filter "status=running" | grep -q server; then
        # 检查HTTP健康检查
        if curl -f http://localhost:8080/health > /dev/null 2>&1; then
            log_success "服务健康检查通过"
        else
            log_warning "服务启动中，请稍后再试"
        fi
    else
        log_error "服务未正常启动"
        docker compose logs server
    fi
}

# 显示服务状态
show_status() {
    echo ""
    log_info "服务状态："
    docker compose ps
    echo ""
    log_info "服务日志（最后10行）："
    docker compose logs --tail=10
}

# 停止服务
stop_services() {
    log_info "停止所有服务..."
    docker compose down
    log_success "服务已停止"
}

# 清理所有数据
cleanup_all() {
    read -p "确定要清理所有数据吗？这将删除数据库数据和所有容器。(y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        log_warning "清理所有数据..."
        docker compose down -v --remove-orphans
        docker system prune -f
        log_success "清理完成"
    else
        log_info "取消清理操作"
    fi
}

# 备份数据
backup_data() {
    log_info "备份数据..."
    
    BACKUP_DIR="backup/$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$BACKUP_DIR"
    
    # 备份MySQL
    docker compose exec mysql mysqldump -u root -p$MYSQL_ROOT_PASSWORD devops_asset > "$BACKUP_DIR/mysql_backup.sql"
    
    # 备份Redis
    docker compose exec redis redis-cli --rdb - > "$BACKUP_DIR/redis_backup.rdb"
    
    log_success "数据备份完成: $BACKUP_DIR"
}

# 恢复数据
restore_data() {
    if [ -z "$1" ]; then
        log_error "请指定备份目录"
        exit 1
    fi
    
    BACKUP_DIR="$1"
    
    if [ ! -d "$BACKUP_DIR" ]; then
        log_error "备份目录不存在: $BACKUP_DIR"
        exit 1
    fi
    
    log_info "恢复数据从: $BACKUP_DIR"
    
    # 恢复MySQL
    if [ -f "$BACKUP_DIR/mysql_backup.sql" ]; then
        docker compose exec -T mysql mysql -u root -p$MYSQL_ROOT_PASSWORD devops_asset < "$BACKUP_DIR/mysql_backup.sql"
        log_success "MySQL数据恢复完成"
    fi
    
    # 恢复Redis
    if [ -f "$BACKUP_DIR/redis_backup.rdb" ]; then
        docker cp "$BACKUP_DIR/redis_backup.rdb" devops-asset-redis:/data/dump.rdb
        docker compose restart redis
        log_success "Redis数据恢复完成"
    fi
}

# 显示帮助信息
show_help() {
    cat << EOF
DevOps Asset Management System 部署脚本

使用方法:
    $0 [选项] [命令]

命令:
    production      部署生产环境
    development     部署开发环境  
    status          显示服务状态
    health          执行健康检查
    stop            停止所有服务
    restart         重启所有服务
    logs            查看服务日志
    backup          备份数据
    restore [dir]   恢复数据
    cleanup         清理所有数据
    help            显示帮助信息

选项:
    --auto-config   自动生成安全配置
    --skip-check    跳过依赖检查

示例:
    $0 production --auto-config    # 自动配置并部署生产环境
    $0 development                 # 部署开发环境
    $0 status                      # 查看服务状态
    $0 backup                      # 备份数据
    $0 restore backup/20240101_120000  # 恢复指定备份

EOF
}

# 主函数
main() {
    # 解析参数
    AUTO_CONFIG=false
    SKIP_CHECK=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --auto-config)
                AUTO_CONFIG=true
                shift
                ;;
            --skip-check)
                SKIP_CHECK=true
                shift
                ;;
            production|development|status|health|stop|restart|logs|backup|restore|cleanup|help)
                COMMAND=$1
                shift
                break
                ;;
            *)
                log_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 如果没有指定命令，显示帮助
    if [ -z "$COMMAND" ]; then
        show_help
        exit 0
    fi
    
    # 检查依赖（除非跳过）
    if [ "$SKIP_CHECK" != "true" ]; then
        check_dependencies
    fi
    
    # 执行命令
    case $COMMAND in
        production)
            setup_environment
            create_directories
            auto_configure_env
            deploy_production
            health_check
            show_status
            ;;
        development)
            setup_environment
            create_directories
            auto_configure_env
            deploy_development
            health_check
            show_status
            ;;
        status)
            show_status
            ;;
        health)
            health_check
            ;;
        stop)
            stop_services
            ;;
        restart)
            stop_services
            sleep 5
            if [ -f .env ] && grep -q "APP_MODE=prod" .env; then
                deploy_production
            else
                deploy_development
            fi
            ;;
        logs)
            docker compose logs -f
            ;;
        backup)
            backup_data
            ;;
        restore)
            restore_data "$1"
            ;;
        cleanup)
            cleanup_all
            ;;
        help)
            show_help
            ;;
        *)
            log_error "未知命令: $COMMAND"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@" 
