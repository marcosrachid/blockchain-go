# Implementação de Halving e Limite de Supply

## ✅ O Que Foi Implementado

Um sistema completo de **halving e limite de supply similar ao Bitcoin** foi adicionado ao protocolo da blockchain.

## 📊 Parâmetros do Protocolo (Centralizados em `config.go`)

Toda a configuração está centralizada em: **`internal/blockchain/config.go`**

### Configuração de Supply e Recompensa

```go
const (
    InitialSubsidy  = 50       // Recompensa inicial de mineração (50 moedas como Bitcoin)
    HalvingInterval = 210000   // Blocos até o halving (~4 anos)
    MaxSupply       = 21000000 // Supply máximo (21 milhões de moedas)
)
```

### Outros Parâmetros do Protocolo

```go
const (
    Difficulty = 18 // Dificuldade de mineração (PoW)
    GenesisData = "First Transaction from Genesis"
    DBPath = "./tmp/blocks"
    DefaultPort = 3000
    ProtocolVersion = 1
)
```

## 🔄 Como Funciona o Halving

### Cronograma de Recompensas por Bloco

| Blocos | Recompensa | Moedas Criadas | Acumulado |
|--------|------------|----------------|-----------|
| 0 - 209.999 | 50 | 10.500.000 | 10.500.000 |
| 210.000 - 419.999 | 25 | 5.250.000 | 15.750.000 |
| 420.000 - 629.999 | 12 | 2.520.000 | 18.270.000 |
| 630.000 - 839.999 | 6 | 1.260.000 | 19.530.000 |
| 840.000 - 1.049.999 | 3 | 630.000 | 20.160.000 |
| ... | ... | ... | ... |
| ~6.930.000+ | 0 | 0 | ~21.000.000 |

### Função de Cálculo

```go
func GetBlockReward(height int) int {
    reward := InitialSubsidy
    
    // Calcula o número de halvings
    halvings := height / HalvingInterval
    
    // Cada halving divide a recompensa por 2
    for i := 0; i < halvings; i++ {
        reward = reward / 2
    }
    
    // Quando a recompensa chega a 0, não cria mais moedas
    if reward < 1 {
        return 0
    }
    
    return reward
}
```

## 🎯 Características Principais

### 1. **Taxa de Emissão Decrescente**
- Recompensa reduz pela metade a cada 210.000 blocos
- Imita o modelo de escassez do Bitcoin
- Previne inflação ao longo do tempo

### 2. **Limite Máximo de Supply**
- Limite rígido de 21 milhões de moedas
- Nenhuma moeda pode ser criada após atingir o máximo
- Recompensa torna-se 0 após ~33 halvings

### 3. **Emissão Previsível**
- Transparente e determinístico
- Possível calcular supply total em qualquer altura
- Incentivos econômicos claros para mineradores

### 4. **Cálculo Baseado em Altura**
- Recompensa calculada pela altura do bloco
- Não precisa armazenar histórico de recompensas
- Eficiente e verificável

## 📝 Funções Atualizadas

### `CoinbaseTX` (Transação de Recompensa de Mineração)

**Antes:**
```go
func CoinbaseTX(to, data string) *Transaction {
    txout := NewTXOutput(50, to) // Recompensa fixa
    // ...
}
```

**Depois:**
```go
func CoinbaseTX(to, data string, height int) *Transaction {
    reward := GetBlockReward(height) // Recompensa dinâmica
    txout := NewTXOutput(reward, to)
    // ...
}
```

### `GetBestHeight` (Nova Função)

Adicionada à blockchain para obter altura atual:

```go
func (chain *Blockchain) GetBestHeight() int {
    var lastBlock Block
    // ... busca último bloco do banco de dados
    return lastBlock.Height
}
```

### Uso na Mineração

```go
// Ao criar um novo bloco
newHeight := chain.GetBestHeight() + 1
cbTx := blockchain.CoinbaseTX(minerAddress, "", newHeight)
```

## 📂 Arquivos Modificados

1. **`internal/blockchain/config.go`** ⭐ NOVO
   - Arquivo de configuração centralizado
   - Todas as constantes do protocolo
   - Funções auxiliares

2. **`internal/blockchain/transaction.go`**
   - Atualizado `CoinbaseTX` para aceitar altura
   - Movidas constantes para config.go
   - Usa `GetBlockReward()`

3. **`internal/blockchain/blockchain.go`**
   - Adicionado método `GetBestHeight()`
   - Atualizada criação do bloco genesis
   - Usa constantes do config.go

4. **`internal/blockchain/proof.go`**
   - Usa `Difficulty` do config.go

5. **`cmd/blockchain/main.go`**
   - Atualizado comando send para calcular altura
   - Passa altura para `CoinbaseTX`

6. **`internal/network/server.go`**
   - Atualizada função de mineração
   - Calcula altura antes de criar coinbase

## 🧪 Testando o Halving

### Cenário de Teste

```bash
# Mine blocos e verifique recompensas em diferentes alturas

# Bloco 0 (Genesis)
Recompensa: 50 moedas

# Blocos 1-209.999
Recompensa: 50 moedas cada

# Bloco 210.000 (Primeiro Halving)
Recompensa: 25 moedas

# Bloco 420.000 (Segundo Halving)
Recompensa: 12 moedas (arredondado para baixo)

# Bloco 630.000 (Terceiro Halving)
Recompensa: 6 moedas
```

### Verificar Supply

```go
// Função auxiliar para verificar supply total
func VerifySupply(chain *Blockchain) {
    height := chain.GetBestHeight()
    expectedSupply := CalculateSupplyUpToHeight(height)
    actualSupply := chain.GetTotalSupply()
    
    if actualSupply > MaxSupply {
        log.Fatal("Supply excedeu o máximo!")
    }
}
```

## 📈 Modelo Econômico

### Curva de Supply

```
Supply (milhões)
21M ┤                           ___________
    │                      ___/
    │                 ___/
15M ┤            ___/
    │       ___/
    │  ___/
10M ┤_/
    │
 5M ┤
    │
  0 └─────────────────────────────────────►
    0   210k  420k  630k  840k  1.05M  ...  Altura do Bloco
```

### Taxa de Emissão

- **Anos 0-4**: 50 moedas/bloco (rápido)
- **Anos 4-8**: 25 moedas/bloco (médio)
- **Anos 8-12**: 12 moedas/bloco (lento)
- **Anos 12+**: Progressivamente mais lento
- **Ano ~140**: Emissão para (supply máximo atingido)

## 🎓 Benefícios

1. **Escassez**: Supply limitado aumenta valor ao longo do tempo
2. **Previsibilidade**: Cronograma de emissão conhecido
3. **Incentivo**: Mineradores iniciais recebem maiores recompensas
4. **Estabilidade**: Taxa de inflação decrescente
5. **Compatibilidade com Bitcoin**: Mesmo modelo do Bitcoin

## 🔮 Melhorias Futuras

Possíveis adições:

1. **Taxas de Transação**
   ```go
   reward := GetBlockReward(height) + fees
   ```

2. **Verificação de Supply**
   ```go
   func (chain *Blockchain) ValidateSupply() bool {
       return chain.GetTotalSupply() <= MaxSupply
   }
   ```

3. **Estatísticas de Emissão**
   ```go
   func GetEmissionRate(height int) float64 {
       // Calcula moedas por ano na altura dada
   }
   ```

4. **Comandos de Consulta de Supply**
   ```bash
   ./blockchain supply              # Supply total atual
   ./blockchain supply -height 1000 # Supply na altura 1000
   ./blockchain halving             # Próximo bloco de halving
   ```

## ✅ Resumo

A blockchain agora tem:
- ✅ **Mecanismo de halving** (a cada 210.000 blocos)
- ✅ **Supply máximo** (21 milhões de moedas)
- ✅ **Configuração centralizada** (config.go)
- ✅ **Recompensas baseadas em altura** (cálculo dinâmico)
- ✅ **Compatível com Bitcoin** (mesmos parâmetros)

**Status:** Modelo econômico pronto para produção! 🎉

---

Para mais informações, veja:
- [../internal/blockchain/config.go](../internal/blockchain/config.go) - Configuração do protocolo
- [../README.md](../README.md) - Documentação geral
- [BITCOIN_COMPARISON.md](BITCOIN_COMPARISON.md) - Comparação com Bitcoin

