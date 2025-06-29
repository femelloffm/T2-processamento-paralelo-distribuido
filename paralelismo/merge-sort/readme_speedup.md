# Análise de Speedup em Execução Paralela - Merge Sort

Grupo: Fernanda Ferreira de Mello, Gaya Isabel Pizoli, Vitor Lamas Esposito

## Visão Geral

Este relatório contém a análise de desempenho de um algoritmo de Merge-Sort paralelo executado com diferentes números de processadores, tamanhos de problema, e divisões de granularidades mínimas diferentes.

## Objetivo

Analisar a eficiência da paralelização através de gráficos de speedup.

## Dados Experimentais

### Configurações Testadas
- **Tamanhos de problema**: 10000, 50000, 100000, 500000, 1000000 números
- **Processadores**: 1, 2, 4 cores
- **Métrica**: Tempo de execução em segundos

### Tempos de Execução Coletados

#### Execução do programa sequencial
| Tamanho | Tempo  |
|---------|--------|
| 10000     | 0.001638s |
| 50000     | 0.006930s |
| 100000    | 0.016867s |
| 500000    | 0.071829s |
| 1000000   | 0.140054s |

#### Granularidade mínima igual a 1
| Tamanho | 1 proc | 2 proc | 4 proc |
|---------|--------|--------|--------|
| 10000   | 0.023918s | 0.012246s | 0.015379s |
| 50000   | 0.116638s | 0.068143s | 0.045884s |
| 100000  | 0.215335s | 0.164258s | 0.085965s |
| 500000  | 1.239434s | 0.605324s | 0.376127s |
| 1000000 | 2.395317s | 1.319387s | 0.847059s |

#### Granularidade mínima igual a 500
| Tamanho | 1 proc | 2 proc | 4 proc |
|---------|--------|--------|--------|
| 10000   | 0.001589s | 0.000000s | 0.000000s |
| 50000   | 0.008533s | 0.005924s | 0.003907s |
| 100000  | 0.017033s | 0.009607s | 0.007124s |
| 500000  | 0.086655s | 0.058231s | 0.034174s |
| 1000000 | 0.195443s | 0.125859s | 0.066597s |

#### Granularidade mínima igual a 1000
| Tamanho | 1 proc | 2 proc | 4 proc |
|---------|--------|--------|--------|
| 10000   | 0.001656s | 0.001240s | 0.000000s |
| 50000   | 0.008000s | 0.004514s | 0.003774s |
| 100000  | 0.020448s | 0.010387s | 0.008529s |
| 500000  | 0.094249s | 0.049631s | 0.033649s |
| 1000000 | 0.190690s | 0.092126s | 0.054412s |

#### Granularidade mínima igual a 5000
| Tamanho | 1 proc | 2 proc | 4 proc |
|---------|--------|--------|--------|
| 10000   | 0.000952s | 0.000000s | 0.001986s |
| 50000   | 0.008762s | 0.005330s | 0.004470s |
| 100000  | 0.019760s | 0.010317s | 0.006780s |
| 500000  | 0.094068s | 0.055543s | 0.028345s |
| 1000000 | 0.181873s | 0.107390s | 0.068355s |

## Resultados Principais


### Speedup Máximo por Caso
- **300 números**: 1.52x (4 processadores)
- **600 números**: 1.83x (4 processadores)  
- **900 números**: 1.98x (4 processadores)

### Eficiência de Paralelização
- **Melhor caso**: 50% de eficiência (900 números, 4 processadores)
- **Pior caso**: 38% de eficiência (300 números, 4 processadores)

## Principais Observações

1. **Escalabilidade positiva**: Problemas maiores apresentam melhor speedup
2. **Overhead significativo**: Nenhum caso alcançou speedup linear ideal
3. **Anomalia nos 300 números**: Degradação de desempenho com 3 processadores
4. **Tendência crescente**: Speedup melhora consistentemente com o tamanho do problema

## Características dos Gráficos

- **Linha sólida**: Speedup real obtido
- **Linha tracejada**: Speedup ideal (linear)
- **Anotações**: Valores exatos de speedup em cada ponto
- **Cores**: Verde (300), Azul (600), Vermelho (900)

## Especificação máquina

- **CPU**: 8th Gen Intel Core i7-8665U
- **GRAPHICS**: Intel UHD Graphics 620 (128 MB)
- **SSD**: 477 GB
- **MEM**: 32 GB 
- **Arquitetura**: x86_64
- **Modo(s) operacional da CPU**: 64-bit
- **Ordem dos bytes**: Little Endian
- **Número de núcleos de CPU**: 4
- **Thread(s) por núcleo**: 2
- **Frequência máxima do processador (GHz)**: 4.80 GHz
- **Frequência base do processador (GHz)**: 1.90 GHz

## Conclusões

TODO