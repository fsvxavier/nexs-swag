#!/bin/bash

echo "=== Instalador do nexs-swag ==="
echo ""

# Verificar se Go está instalado
if ! command -v go &> /dev/null; then
    echo "❌ Go não está instalado"
    echo "   Instale Go 1.24+ em: https://golang.org/dl/"
    exit 1
fi

echo "✓ Go instalado: $(go version)"
echo ""

# Instalar
echo "Instalando nexs-swag..."
go install ./cmd/nexs-swag

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Instalação concluída!"
    echo ""
    echo "Binário instalado em: $(go env GOPATH)/bin/nexs-swag"
    echo ""
    
    # Verificar se está no PATH
    if command -v nexs-swag &> /dev/null; then
        echo "✓ nexs-swag está no PATH"
        echo ""
        echo "Teste com: nexs-swag --version"
        nexs-swag --version
    else
        echo "⚠️  nexs-swag não está no PATH"
        echo ""
        echo "Adicione ao seu ~/.bashrc ou ~/.zshrc:"
        echo "  export PATH=\$PATH:\$(go env GOPATH)/bin"
        echo ""
        echo "Ou execute:"
        echo "  export PATH=\$PATH:$(go env GOPATH)/bin"
        echo "  nexs-swag --version"
    fi
else
    echo ""
    echo "❌ Erro na instalação"
    exit 1
fi
