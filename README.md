# planwiadomienia
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/open-source.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/you-didnt-ask-for-this.svg)](https://forthebadge.com)

## O projekcie
webhook na discordzie cie prześladuje kiedy są lekcje 

## Po co ten szajs?
bo nie lubię siedzieć przy komputerze podczas przerwy i leżę w łóżku 

# Użycie
Zmieniasz nazwę pliku `config.example.json` na `config.json`, poprawiasz wszystkie wartości (godzina pierwszej lekcji, długość przerw, treść wiadomości i link do webhooka z Discorda), wpisujesz tam swój plan lekcji (1 - poniedziałek, 2 - wtorek, ..., 5 - piątek; 6 - sobota, 0 - niedziela (współczuję jeżeli ktoś ma lekcje w te dwa dni)) zgodnie z przykładem; jeżeli masz na późniejszą godzinę niż `godzina pierwszej lekcji`, masz okienko lub dwie lekcje pod rząd bez przerwy (bo was nauczyciel nie szanuje) to wpisz `null` (nie `"null"`!) zamiast nazwy lekcji.

Później `go build .` (ewentualnie wcześniej zmieniasz `GOOS` i `GOARCH` w zależności od tego gdzie chcesz to uruchomić) i uruchamiasz plik, który komenda wypluje ~~chyba nie powinnam tego tłumaczyć no ale za późno~~

btw [fratik](https://github.com/fratik) to niewoln- znaczy pomógł mi 
<!-- golang to gówno i strasznie nienawidze czasu -->
<!-- to używaj js!!!!! btw juby cute af -->
