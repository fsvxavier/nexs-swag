#!/bin/bash

echo "=== Exemplo 20: Fmt Command ==="
echo ""

# Criar cópia para demonstrar formatting
cp main.go main_unformatted.go

# Desformatar propositalmente
cat > main_unformatted.go << 'EOF'
package main

import (
	"encoding/json"
	"net/http"
)

// @title Fmt Command Demo API
// @version 1.0

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
}

// GetProduct returns a product
// @Summary Get product
// @Description Get product by ID
// @Tags products
// @Param id path int true "Product ID"
// @Success 200 {object} Product
// @Router /products/{id} [get]
func GetProduct(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Product{ID: 1, Name: "Test"})
}

func main() {
	http.HandleFunc("/api/products/", GetProduct)
	http.ListenAndServe(":8080", nil)
}
EOF

echo "Antes do fmt:"
echo "-------------"
grep -A 7 "// GetProduct" main_unformatted.go

echo ""
echo "Executando fmt..."
../../nexs-swag fmt --dir .

echo ""
echo "Depois do fmt:"
echo "--------------"
grep -A 10 "// GetProduct" main.go

echo ""
echo "✓ Annotations formatadas com alinhamento!"
echo ""
echo "O comando fmt:"
echo "  • Alinha annotations"
echo "  • Organiza espaçamento"
echo "  • Melhora legibilidade"
echo "  • Não altera lógica"

# Limpar
rm -f main_unformatted.go
