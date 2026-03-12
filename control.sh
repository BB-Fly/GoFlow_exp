#!/bin/bash

# 控制GoFlow服务的启动和停止

# 定义服务名称
SERVICE_NAME="GoFlow_exp"

# 检查Redis命令是否存在
check_redis() {
    echo "Checking if Redis is installed..."
    if ! command -v redis-server &> /dev/null; then
        echo "Redis is not installed. Please install Redis first."
        echo "You can install Redis using Homebrew: brew install redis"
        return 1
    fi
    return 0
}

# 启动Redis服务
start_redis() {
    echo "Checking Redis service..."
    if ! pgrep -x "redis-server" > /dev/null; then
        echo "Redis is not running. Starting Redis..."
        redis-server --daemonize yes
        if [ $? -eq 0 ]; then
            echo "Redis started successfully"
        else
            echo "Failed to start Redis. Please check your Redis installation."
            return 1
        fi
    else
        echo "Redis is already running"
    fi
    return 0
}

# 启动GoFlow服务
start_service() {
    # 检查Redis
    check_redis || return 1
    
    # 启动Redis
    start_redis || return 1
    
    # 等待Redis服务完全启动
    sleep 2
    
    # 检查可执行文件是否存在
    if [ -f "./GoFlow_exp" ]; then
        EXECUTABLE="./GoFlow_exp"
    elif [ -f "./build/GoFlow_exp" ]; then
        EXECUTABLE="./build/GoFlow_exp"
    else
        echo "GoFlow_exp executable not found! Please run build.sh first."
        return 1
    fi
    
    # 检查服务是否已经在运行
    if pgrep -x "$SERVICE_NAME" > /dev/null; then
        echo "GoFlow service is already running"
        return 0
    fi
    
    # 启动服务
    echo "Starting GoFlow service..."
    $EXECUTABLE &
    
    # 保存进程ID
    echo $! > .goflow.pid
    echo "GoFlow service started with PID $(cat .goflow.pid)"
    return 0
}

# 停止GoFlow服务
stop_service() {
    # 检查服务是否在运行
    if [ -f ".goflow.pid" ]; then
        PID=$(cat .goflow.pid)
        if kill -0 $PID 2>/dev/null; then
            echo "Stopping GoFlow service with PID $PID..."
            kill $PID
            if [ $? -eq 0 ]; then
                echo "GoFlow service stopped successfully"
                rm .goflow.pid
            else
                echo "Failed to stop GoFlow service"
                return 1
            fi
        else
            echo "GoFlow service is not running, but PID file exists. Cleaning up..."
            rm .goflow.pid
        fi
    else
        # 尝试通过进程名停止
        if pgrep -x "$SERVICE_NAME" > /dev/null; then
            echo "Stopping GoFlow service..."
            pkill -x "$SERVICE_NAME"
            if [ $? -eq 0 ]; then
                echo "GoFlow service stopped successfully"
            else
                echo "Failed to stop GoFlow service"
                return 1
            fi
        else
            echo "GoFlow service is not running"
        fi
    fi
    return 0
}

# 检查服务状态
status_service() {
    if [ -f ".goflow.pid" ]; then
        PID=$(cat .goflow.pid)
        if kill -0 $PID 2>/dev/null; then
            echo "GoFlow service is running with PID $PID"
            return 0
        else
            echo "GoFlow service is not running, but PID file exists. Cleaning up..."
            rm .goflow.pid
            return 1
        fi
    else
        if pgrep -x "$SERVICE_NAME" > /dev/null; then
            echo "GoFlow service is running"
            return 0
        else
            echo "GoFlow service is not running"
            return 1
        fi
    fi
}

# 主函数
main() {
    case "$1" in
        start)
            start_service
            ;;
        stop)
            stop_service
            ;;
        status)
            status_service
            ;;
        restart)
            stop_service
            start_service
            ;;
        *)
            echo "Usage: $0 {start|stop|status|restart}"
            echo "  start   - Start the GoFlow service"
            echo "  stop    - Stop the GoFlow service"
            echo "  status  - Check the status of the GoFlow service"
            echo "  restart - Restart the GoFlow service"
            return 1
            ;;
    esac
}

# 执行主函数
main "$@"
