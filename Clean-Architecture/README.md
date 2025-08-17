## Objetivo
Preparar o esqueleto do projeto **orders** para evoluir sem débito técnico:
- Layout de pastas
- `go.mod`
- `Makefile` com alvos básicos
- `.env.example`
- Binário mínimo que compila

## Requisitos
- Go 1.21+

## Como usar

```bash
# 1) ajuste o módulo no go.mod (opcional)
```
sed -i 's|github.com/seuuser/orders|<seu/modulo>|' go.mod
```

# 2) instalar deps vazias e organizar (não há deps por enquanto)
```
make tidy
```

# 3) compilar
```
make build
```

# 4) executar
```
make run
```

Saída esperada:
```nginx
orders — Passo 0 OK. Pronto para evoluir para domínio e casos de uso.
```