# Generate code

ModelHelper:

Hva skal skje når man genererer kode.
Dette skal gjøres i generate command

## Før

- Hent inn config
- Lag Context
- Last inn alle kildekoblinger
- Finn korrekt kilde basert på input
- Last inn alle språk definisjoner
- Last inn alle maler
- Finn mal- blokker

### For alle valgte maler

- bestem hvilken mal (basert på maler,
- bestem "input"

#### for hver valgte kilde

- lag korrekt input
- Finn korrekt filbane
- Eksisterer filen fra før?
- Hvis eksisterer, er det noen forskjeller?
- Hvis forskjell - hvor er forskjellen
- Skriv til resultatfil

## Etter

- Eksporter hvis valgt
- Vis på skjerm hvis valgt
- Vis statistikk

## På sikt

Skriv til database i prosjekt på struktur
Lag kode kun på diff


## Hent inn config

Dette kan gjøres på flere måter.

- Hvis -c | --config er angitt så skal denne brukes
- Last inn fra standard lokasjon

## Lag Context

Contexten lages på bakgrunn av hvilken kommando som blir kjørt.
Denne henter inn config basert på retningslinjer

Bør dette være en egen interface? på den måten kan jeg operere med flere ulike contexter basert på hva behovet er?

## Last inn alle kildekoblinger

