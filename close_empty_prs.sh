#!/bin/bash

# Script para fechar Pull Requests vazios
# Uso: ./close_empty_prs.sh SEU_TOKEN_AQUI

if [ -z "$1" ]; then
    echo "Erro: Token do GitHub não fornecido"
    echo "Uso: ./close_empty_prs.sh SEU_TOKEN_AQUI"
    echo ""
    echo "Para obter um token:"
    echo "1. Vá para https://github.com/settings/tokens"
    echo "2. Clique em 'Generate new token (classic)'"
    echo "3. Selecione o escopo 'repo'"
    echo "4. Copie o token e use neste script"
    exit 1
fi

TOKEN="$1"
REPO="pamyydev/eco-link"

echo "Fechando Pull Requests vazios..."

# Lista dos PRs vazios para fechar
PRS_TO_CLOSE=(6 7 8 11)
PR_NAMES=("feature/improve-greenweb-adapter" "feature/add-retry-system" "feature/add-structured-logging" "feature/improve-websitecarbon-adapter")

for i in "${!PRS_TO_CLOSE[@]}"; do
    PR_NUMBER="${PRS_TO_CLOSE[$i]}"
    PR_NAME="${PR_NAMES[$i]}"
    
    echo "Fechando PR #$PR_NUMBER ($PR_NAME)..."
    
    RESPONSE=$(curl -s -X PATCH \
        -H "Authorization: token $TOKEN" \
        -H "Accept: application/vnd.github.v3+json" \
        "https://api.github.com/repos/$REPO/pulls/$PR_NUMBER" \
        -d '{"state": "closed"}')
    
    if echo "$RESPONSE" | grep -q '"state":"closed"'; then
        echo "PR #$PR_NUMBER fechado com sucesso!"
    else
        echo "Erro ao fechar PR #$PR_NUMBER"
        echo "Resposta: $RESPONSE"
    fi
    echo ""
done

echo "Processo concluído!"
echo ""
echo "Pull Requests que permanecem abertos:"
echo "- PR #9 - feature/add-error-tests (2 arquivos alterados)"
echo "- PR #10 - feature/improve-websitecarbon-api (1 arquivo alterado)"
echo ""
echo "Próximos passos:"
echo "1. Revise os PRs #9 e #10"
echo "2. Faça merge quando estiver satisfeito"
echo "3. Teste as funcionalidades após o merge"
