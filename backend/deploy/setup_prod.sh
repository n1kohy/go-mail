#!/bin/bash
# setup_prod.sh - 生产环境配置切换脚本

# 设置颜色
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}开始修改配置为生产环境 (Docker)...${NC}"

# 此脚本需在 backend 根目录下运行（或确保能访问到 services/）
if [ ! -d "services" ]; then
    echo "错误：请在 backend 根目录下运行此脚本"
    exit 1
fi

# 1. 替换数据库和Redis地址 (localhost -> host.docker.internal)
# 适用于所有 yaml 文件
echo "正在更新 DB/Redis 地址..."
find services -name "*.yaml" -type f -exec sed -i 's/localhost:5432/host.docker.internal:5432/g' {} +
find services -name "*.yaml" -type f -exec sed -i 's/localhost:6379/host.docker.internal:6379/g' {} +

# 2. 替换 RPC 服务地址 (localhost:99xx -> servicename:99xx)
# 利用 Docker Compose 的服务发现
echo "正在更新 RPC 服务地址..."

# 定义端口到服务名的映射
declare -A rpc_map=(
    ["9901"]="user-rpc"
    ["9902"]="product-rpc"
    ["9903"]="inventory-rpc"
    ["9904"]="cart-rpc"
    ["9905"]="promotion-rpc"
    ["9906"]="order-rpc"
    ["9907"]="payment-rpc"
    ["9908"]="logistics-rpc"
)

# 遍历所有 yaml 文件进行替换
for port in "${!rpc_map[@]}"; do
    service_name=${rpc_map[$port]}
    echo "  映射端口 $port -> $service_name"
    find services -name "*.yaml" -type f -exec sed -i "s/localhost:$port/$service_name:$port/g" {} +
done

echo -e "${GREEN}配置修改完成！${NC}"
echo "请执行 docker-compose up -d --build 开始部署"
