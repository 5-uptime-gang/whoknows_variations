# Monitorering refleksioner

Da vi begyndte at indsamle metrikker fra vores virtuelle maskine, fulgte vi de klassiske systemmålinger - CPU-belastning, diskforbrug og RAM-usage. Det gav dog hurtigt anledning til en vigtig erkendelse: Det skaber ingen reel værdi at monitorere komponenter, der allerede fungerer stabilt. Meningsfuld monitorering skal føre til opdagelser, der enten forebygger fejl eller gør os i stand til at optimere.
Et konkret eksempel opstod, da RAM-monitoreringen viste et konstant forbrug omkring 70%. Vi har tidligere oplevet, at serveren crasher ved 100% RAM-forbrug, og metrikkerne gjorde det tydeligt, at vi var tættere på grænsen end forventet. Det førte til en ændring i vores tilgang: Vi vil gerne analysere, hvilke processer og handlinger der driver RAM-forbruget - eksempelvis deployment-opstart og API-trafik. Denne indsigt giver os derfor mulighed for at forebygge nedbrud og planlægge optimeringer baseret på faktiske data.


## Antal resultater pr. søgning
Dette giver indblik i, om vores søgninger er relevante og dækkende. For mange eller for få resultater peger på problemer i datagrundlag, søgelogik eller filtrering, som kan påvirke brugeroplevelsen direkte.

## Trending søgninger
Ved at se hvilke søgetermer der trender, får vi viden om, hvad brugerne faktisk efterspørger. Det hjælper os med at prioritere udvikling og forbedre indhold og funktioner, så de matcher brugernes behov.

## Antal oprettede brugere
Oprettelsesraten viser, om vores platform tiltrækker nye brugere og om onboarding fungerer. Ændringer i tallet kan afsløre fejl, barrierer eller behov for forbedringer i registreringsflowet.


## Login-succesrate
Ved at måle hvor mange der logger ind succesfuldt, kan vi identificere problemer i loginflowet. Mange fejlslagne forsøg tyder på tekniske fejl, dårlig brugerforståelse eller behov for tydeligere feedback.


## Retry-adfærd ved loginfejl
At tracke om brugere forsøger igen efter fejl, fortæller om oplevet brugervenlighed. Stopper de efter første fejl, kan det indikere frustration, mens gentagne forsøg kan pege på misforståelser eller dårlig guidance.


## Trafikmønstre / peak-usage
At vide hvornår trafikken topper gør det lettere at planlægge ressourcer og skalering. Det hjælper også med at forstå brugeradfærd på tværs af døgnet.

## Ekstra
Disse områder kan også være dejlige at monitorere, men det er kun hvis tiden tillader det:

## Loading time
Indlæsningshastighed påvirker direkte, om brugerne bliver eller forlader siden. Længere load times afslører flaskehalse, som kan optimeres for både performance og brugeroplevelse.

## Tid brugerne spenderer på sitet
Dette viser, om indholdet engagerer. Kort tid på siden kan indikere, at brugerne ikke finder det, de søger, eller at navigationen ikke fungerer optimalt.

## Brugere der forlader uden at interagere
At måle hvor mange der forlader uden klik, afslører problemer med relevans, layout eller førstehåndsindtryk. Det er en stærk indikator på, om siden fanger brugerens opmærksomhed.


## Browservalg (accessibility)
Viden om hvilke browsere brugerne anvender hjælper os med at prioritere kompatibilitet og fejlrettelser. Det sikrer, at flest muligt får en stabil og tilgængelig oplevelse.

## Første element brugerne scroller/klikker mod
Dette viser, hvad der faktisk tiltrækker brugerens opmærksomhed på siden. Det hjælper med at placere vigtigt indhold rigtigt og optimere layout.


## Kvaliteten af søgeresultater (klikrate på resultater)
Hvis brugerne klikker på de øverste resultater, er rangeringen god. Hvis ikke, tyder det på, at søgelogikken skal forbedres. Det giver et direkte mål for, hvor brugbare vores søgeresultater er.


## Pull Requests – behandlingstid
At måle hvor længe PR’er ligger åbne, giver et klart billede af flowet i udviklingsprocessen. Lange behandlingstider kan pege på flaskehalse, manglende prioritering eller behov for bedre review-strukturer.

## Issues – tid før løsning
Ved at tracke hvor lang tid issues er åbne, kan vi identificere mønstre i, hvilke typer problemer der tager længst tid, og hvor arbejdet eventuelt stopper op. Det hjælper med at effektivisere planlægning og ressourcefordeling.

## Gamification af udviklingsprocessen
Ved at visualisere og belønne hurtige, konsistente eller særligt vigtige bidrag (fx hurtige PR-reviews eller løsning af ældre issues) kan man motivere teamet og skabe en mere engagerende og målbar udviklingskultur. Det styrker både samarbejde og gennemførselshastighed.

### Database–query performance
Langsomme queries påvirker hele brugeroplevelsen. Overvågning af query-tider og låsninger gør det muligt at optimere SQL og indeks før problemer eskalerer.
