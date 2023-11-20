FROM --platform=linux/arm64/v8 alpine:latest
RUN apt update && apt install -y ffmpeg\ && \ apt-get clean && \ rm -rf /var/lib/apt/lists/*
COPY mainsendergo /
ENTRYPOINT ["/mainsendergo"]
