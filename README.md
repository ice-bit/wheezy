# wheezy ![](https://github.com/ice-bit/wheezy/actions/workflows/wheezy.yml/badge.svg)
wheezy è una webapp scritta in Go per calcolare il codice fiscale e il suo
inverso. L'applicazione è accessibile a [cf.marcocetica.com](https://cf.marcocetica.com).

Il **frontend** è stato realizzato utilizzando Bootstrap, mentre il **backend** è
stato realizzato utilizzando il pacchetto `net/http` di Go. L'intero progetto
segue il pattern MVC.

Questa webapp si avvale di una database SQLite per organizzare i _codici catastali_ e 
i _codici nazionali_ necessari al calcolo del codice fiscale. Tali codici
sono scaricati in maniera automatizzata dal [sito del ministero dell'interno](https://dait.interno.gov.it/territorio-e-autonomie-locali/sut/elenco_codici_comuni.php%22) e 
dal [sito dell'ISTAT](https://www.istat.it/it/archivio/6747). Per 
velocizzarne la lettura, questi codici vengono inoltre salvati in una cache in RAM
utilizzando un database _key-value_(**redis**).

## Deploy
Il metodo di deploy supportato ufficialmente è mediante Docker. Per lanciare la webapp
è dunque sufficiente lanciare il seguente comando:
```sh
$> docker-compose up -d
```

Di default sono utilizzati i seguenti parametri:
-  **WHEEZY_LISTEN_ADDRESS**: `127.0.0.1`;
-  **WHEEZY_LISTEN_PORT**: `9000`;
-  **WHEEZY_REDIS_ADDRESS**: `127.0.0.1`;
-  **WHEEZY_REDIS_PORT**: `6379`.

Nel caso alcune di queste variabili d'ambiente dovessero collidere con le configurazioni
di altri servizi in esecuzione, è possibile andare ad alterarle modificando i file `Dockerfile` e `docker-compose.yml` della root del progetto.

Una volta lanciati i due container(quello della webapp e quello della cache redis), 
è possibile esporre l'applicazione tramite un reverse proxy(in questo esempio, tramite nginx) utilizzando la seguente voce:

```nginx
location / {
    proxy_pass http://127.0.0.1:9000;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

Nel caso si fosse scelto una porta differente, modificare la proprietà `proxy_pass`
in modo appropriato.

## Aggiornamento Databases
Di default la webapp fornisce un database(`codes.db`) contenente due tabelle: una per 
i codici catastali e una per quelli nazionali. Nel caso fosse necessario
aggiornare i codici con dati più recenti(ad esempio a seguito della soppressione di 
un comune), seguire la seguente procedura:

1. Scaricare dal [sito dell'ISTAT](https://www.istat.it/it/archivio/6747)
l'**elenco codici e denominazioni delle unità territoriali estere** in formato `zip` 
alla voce _Elenco codici e denominazioni delle unità territoriali estere_. 
2. Scompattare dall'archivio solo il file in formato `csv`.
3. Il file è in formato `latin1`(ISO-8859), prima di poterlo utilizzare, è necessario convertirlo
in formato UTF8. Per farlo eseguire il comando:
```shell
$> iconv -f ISO-8859-11 tabella.csv -t UTF-8 -o tabella.csv
```
4. Normalizzare il file in formato csv(i.e. sostituire `;` con `,`):
```shell
$> sed -i 's/;/,/g' tabella.csv
```
5. Eliminare le linee vuote utilizzando il seguente comando:
```shell
$> sed -i '/^,/d' tabella.csv
```
6. Convertire il file `csv` in un file `sql` utilizzando lo 
script `codnazioni.py`:
```shell
$> python3 codnazioni.py tabella.csv codnazioni.sql
```
7. Scaricare i codici catastali utilizzando lo script `codcatastali.py`(dipendenze: `beautifulsoup4` e `requests`):
```shell
$> python3 codcatastali.py codcatastali.sql
```
8. Creare il database con il seguente comando:
```shell
$> sqlite3 codes.db < cod{nazioni,catastali}.sql
```
Il file `codes.db`, presente nella root del progetto, verrà inserito all'interno del container durante
la fase di deploy.

## Sviluppo
Per poter modificare, aggiungere o debuggare le funzionalità della webapp, è necessario
aver installato una copia del compilatore di Go e una recente versione del server Redis.
È altresì necessario configurare le seguenti variabili d'ambiente:
-  **WHEEZY_LISTEN_ADDRESS**: `127.0.0.1`;
-  **WHEEZY_LISTEN_PORT**: `9000`;
-  **WHEEZY_REDIS_ADDRESS**: `127.0.0.1`;
-  **WHEEZY_REDIS_PORT**: `6379`.

Fatto ciò, è sufficiente avviare l'applicazione utilizzando il seguente comando; 
le (due)dipendenze verranno installate in maniera automatica:

```sh
$> go run ./...
```

## Unit Tests
La _business logic_ del progetto(cartella `model`) è sottoposta a
test unitari i quali permettono di garantire il corretto funzionamento
dell'applicazione. Per poter invocare il motore di testing, eseguire
il seguente comando:
```sh
$> go test ./... -v
```

## Licenza
[GPLv3](https://choosealicense.com/licenses/gpl-3.0/)