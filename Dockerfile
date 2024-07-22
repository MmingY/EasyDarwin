FROM arm64v8/ubuntu:latest

RUN apt update
RUN apt upgrade -y
# RUN DEBIAN_FRONTEND=noninteractive apt install vim lsof curl wget -y
RUN apt install ffmpeg -y
COPY . .
EXPOSE 554
EXPOSE 10008
# 设置容器启动时执行的命令
# CMD ["./start.sh"]
CMD ["./easydarwin"]