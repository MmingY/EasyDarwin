FROM arm64v8/ubuntu:latest
# FROM arm64v8/ubuntu:20.04
# FROM linuxserver/ffmpeg:latest

# 将清华大学的镜像源添加到 /etc/apt/sources.list 文件中//
# RUN sed -i 's/archive.ubuntu.com/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    sed \
    && rm -rf /var/lib/apt/lists/*

RUN sed -i 's/archive.ubuntu.com/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list


RUN apt-get update
RUN apt-get upgrade -y
# RUN DEBIAN_FRONTEND=noninteractive apt install vim lsof curl wget -y
RUN DEBIAN_FRONTEND=noninteractiv apt install -y ffmpeg
COPY . .
EXPOSE 554
EXPOSE 10008
VOLUME /opt/video
# 设置容器启动时执行的命令
# CMD ["./start.sh"]
CMD ["./easydarwin"]