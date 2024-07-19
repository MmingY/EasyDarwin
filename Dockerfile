FROM arm64v8/ubuntu:latest

RUN apt update
RUN apt upgrade -y
RUN DEBIAN_FRONTEND=noninteractive apt install vim lsof curl wget -y
COPY . .
EXPOSE 554
EXPOSE 10086
EXPOSE 10054
EXPOSE 10010
EXPOSE 10035
# 设置容器启动时执行的命令
CMD ["./start.sh"]