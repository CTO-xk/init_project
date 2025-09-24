#!/bin/bash

# Hardhat 智能合约测试脚本
# 使用方法: ./hardhat_test_commands.sh [选项]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的信息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Node.js 和 npm 是否安装
check_dependencies() {
    if ! command -v node &> /dev/null; then
        print_error "Node.js 未安装，请先安装 Node.js"
        exit 1
    fi
    
    if ! command -v npm &> /dev/null; then
        print_error "npm 未安装，请先安装 npm"
        exit 1
    fi
    
    print_success "Node.js 和 npm 已安装"
    node --version
    npm --version
}

# 安装项目依赖
install_dependencies() {
    print_info "检查并安装项目依赖..."
    
    if [ ! -d "node_modules" ]; then
        print_info "安装项目依赖..."
        npm install
    else
        print_info "项目依赖已安装"
    fi
}

# 编译合约
compile_contracts() {
    print_info "编译智能合约..."
    
    # 使用 Hardhat 编译
    npx hardhat compile
    
    print_success "合约编译完成"
}

# 运行基础测试
basic_tests() {
    print_info "开始运行基础测试..."
    
    # 编译合约
    compile_contracts
    
    # 运行测试
    print_info "运行测试..."
    npx hardhat test
    
    print_success "基础测试完成"
}

# 运行特定测试文件
specific_tests() {
    print_info "开始运行特定测试..."
    
    # 编译合约
    compile_contracts
    
    # 运行原始合约测试
    print_info "测试原始合约..."
    npx hardhat test test/ArithmeticTest.sol
    
    # 运行优化合约测试
    print_info "测试优化合约..."
    npx hardhat test test/ArithmeticOptimizedTest.sol
    
    print_success "特定测试完成"
}

# 运行 Gas 分析测试
gas_tests() {
    print_info "开始运行 Gas 分析测试..."
    
    # 编译合约
    compile_contracts
    
    # 运行测试并显示 Gas 信息
    print_info "运行 Gas 分析测试..."
    npx hardhat test test/GasAnalysisTest.sol
    
    print_success "Gas 分析测试完成"
}

# 运行详细测试
verbose_tests() {
    print_info "开始运行详细测试..."
    
    # 编译合约
    compile_contracts
    
    # 运行测试并显示详细输出
    print_info "运行详细测试..."
    npx hardhat test --verbose
    
    print_success "详细测试完成"
}

# 启动本地网络
start_local_network() {
    print_info "启动本地 Hardhat 网络..."
    
    # 在后台启动网络
    npx hardhat node &
    NETWORK_PID=$!
    
    print_info "本地网络已启动，PID: $NETWORK_PID"
    print_info "网络地址: http://127.0.0.1:8545"
    print_info "使用 Ctrl+C 停止网络"
    
    # 等待网络启动
    sleep 3
    
    # 部署合约
    print_info "部署合约到本地网络..."
    npx hardhat run script/deploy.js --network localhost
    
    # 保持网络运行
    wait $NETWORK_PID
}

# 部署合约
deploy_contracts() {
    print_info "部署合约..."
    
    # 部署到本地网络
    print_info "部署到本地网络..."
    npx hardhat run script/deploy.js --network localhost
    
    print_success "合约部署完成"
}

# 清理和重置
clean_and_reset() {
    print_info "清理项目..."
    
    # 清理编译缓存
    npx hardhat clean
    
    # 清理 node_modules（可选）
    read -p "是否删除 node_modules 目录？(y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_info "删除 node_modules 目录..."
        rm -rf node_modules
        print_info "重新安装依赖..."
        npm install
    fi
    
    print_success "清理完成"
}

# 显示帮助信息
show_help() {
    echo -e "${BLUE}Hardhat 智能合约测试脚本${NC}"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  help          显示此帮助信息"
    echo "  install       安装项目依赖"
    echo "  compile       编译合约"
    echo "  basic         运行基础测试"
    echo "  specific      运行特定测试"
    echo "  gas           运行 Gas 分析测试"
    echo "  verbose       运行详细测试"
    echo "  network       启动本地网络"
    echo "  deploy        部署合约"
    echo "  clean         清理项目"
    echo "  all           运行完整测试套件"
    echo ""
    echo "示例:"
    echo "  $0 basic              # 运行基础测试"
    echo "  $0 specific           # 运行特定测试"
    echo "  $0 gas               # 运行 Gas 分析"
    echo "  $0 all               # 运行完整测试套件"
}

# 运行所有测试
run_all_tests() {
    print_info "开始运行完整测试套件..."
    
    install_dependencies
    echo ""
    compile_contracts
    echo ""
    basic_tests
    echo ""
    specific_tests
    echo ""
    gas_tests
    echo ""
    
    print_success "完整测试套件运行完成！"
}

# 主程序
main() {
    case "${1:-help}" in
        "help"|"-h"|"--help")
            show_help
            ;;
        "install")
            check_dependencies && install_dependencies
            ;;
        "compile")
            check_dependencies && install_dependencies && compile_contracts
            ;;
        "basic")
            check_dependencies && install_dependencies && basic_tests
            ;;
        "specific")
            check_dependencies && install_dependencies && specific_tests
            ;;
        "gas")
            check_dependencies && install_dependencies && gas_tests
            ;;
        "verbose")
            check_dependencies && install_dependencies && verbose_tests
            ;;
        "network")
            check_dependencies && install_dependencies && start_local_network
            ;;
        "deploy")
            check_dependencies && install_dependencies && deploy_contracts
            ;;
        "clean")
            clean_and_reset
            ;;
        "all")
            check_dependencies && run_all_tests
            ;;
        *)
            print_error "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
}

# 运行主程序
main "$@"
