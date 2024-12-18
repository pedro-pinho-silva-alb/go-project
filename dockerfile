FROM debian:bullseye AS builder

# Atualizar e instalar dependências do sistema
RUN apt-get update && apt-get install -y \
    build-essential \
    wget \
    curl \
    git \
    libssl-dev \
    pkg-config

# Baixar e instalar o Go
RUN wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz && \
    rm go1.23.3.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app

# Copiar o restante do código e construir o binário
COPY code/. .

# Instalar dependências do Go
RUN go mod download

RUN go build -o main .

# Preparar a imagem final
FROM debian:bullseye

# Instalar apenas as dependências necessárias para execução
RUN apt-get update && apt-get install -y libssl-dev && apt-get clean

# Copiar o binário da etapa de compilação
COPY --from=builder /app/main .

EXPOSE 8085

CMD ["./main"]
