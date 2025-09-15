#!/bin/bash
# Bluebell 应用一键部署脚本

set -e  # 遇到错误立即退出

echo "🚀 开始部署 Bluebell 应用..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印信息函数
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查是否为 root 用户
if [ "$EUID" -ne 0 ]; then
    error "请使用 root 用户运行此脚本"
    exit 1
fi

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
APP_DIR="$(dirname "$SCRIPT_DIR")"

# 1. 更新系统包
info "更新系统包..."
apt update

# 2. 安装 MySQL
info "检查 MySQL 安装状态..."
if ! command -v mysql &> /dev/null; then
    info "安装 MySQL Server..."
    DEBIAN_FRONTEND=noninteractive apt install -y mysql-server
else
    info "MySQL 已安装，跳过..."
fi

# 3. 安装 Redis
info "检查 Redis 安装状态..."
if ! command -v redis-server &> /dev/null; then
    info "安装 Redis Server..."
    apt install -y redis-server
else
    info "Redis 已安装，跳过..."
fi

# 4. 安装 Supervisor（如果需要）
info "检查 Supervisor 安装状态..."
if ! command -v supervisorctl &> /dev/null; then
    info "安装 Supervisor..."
    apt install -y supervisor
else
    info "Supervisor 已安装，跳过..."
fi

# 5. 配置 MySQL
info "配置 MySQL 数据库..."
# 设置 root 密码为 781129
mysql -e "ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '781129';" 2>/dev/null || true
mysql -u root -p781129 -e "FLUSH PRIVILEGES;" 2>/dev/null || true

# 初始化数据库
info "初始化 Bluebell 数据库..."
mysql -u root -p781129 < "$SCRIPT_DIR/init_database.sql"

# 6. 配置 Redis
info "配置 Redis..."
# 备份原配置文件
cp /etc/redis/redis.conf /etc/redis/redis.conf.bak

# 设置 Redis 密码
sed -i 's/# requirepass foobared/requirepass 781129/' /etc/redis/redis.conf

# 重启 Redis 服务
systemctl restart redis-server
systemctl enable redis-server

# 7. 启动 MySQL 服务
info "启动 MySQL 服务..."
systemctl start mysql
systemctl enable mysql

# 8. 验证服务状态
info "验证服务状态..."
if systemctl is-active --quiet mysql; then
    info "✅ MySQL 服务运行正常"
else
    error "❌ MySQL 服务启动失败"
    exit 1
fi

if systemctl is-active --quiet redis-server; then
    info "✅ Redis 服务运行正常"
else
    error "❌ Redis 服务启动失败"
    exit 1
fi

# 9. 测试数据库连接
info "测试数据库连接..."
if mysql -u root -p781129 -e "USE bluebell; SHOW TABLES;" &> /dev/null; then
    info "✅ MySQL 连接测试成功"
else
    error "❌ MySQL 连接测试失败"
    exit 1
fi

if redis-cli -a 781129 ping &> /dev/null; then
    info "✅ Redis 连接测试成功"
else
    error "❌ Redis 连接测试失败"
    exit 1
fi

# 10. 重启 Supervisor 服务
info "重启 Supervisor 服务..."
if command -v supervisorctl &> /dev/null; then
    supervisorctl restart bluebell || warn "无法重启 bluebell 服务，请检查 supervisor 配置"
    info "✅ 尝试重启 Bluebell 应用"
else
    warn "Supervisor 未正确配置，请手动启动应用"
fi

echo ""
info "🎉 部署完成！"
info "MySQL: localhost:3306, 用户: root, 密码: 781129"
info "Redis: localhost:6379, 密码: 781129"
info "应用目录: $APP_DIR"
info "查看应用日志: supervisorctl tail bluebell"
echo ""