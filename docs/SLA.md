# Service Level Agreement for femprocentoppetid.dk

## Service Scope
Denne aftale dækker brugen af webapplikationen femprocentoppetid.dk. Servicen inkluderer:
Søgemaskine der lever på en webapplikation.

Afgrænsning:
Denne SLA dækker ikke 3. parts services (som GitHubs egen nedetid), da disse paradoksalt nok hjælper os med at nå vores mål om lav oppetid.

## Performance Metrics
I overensstemmelse med servicens navn og brand-identitet, måles performance ud fra evnen til at holde systemet nede størstedelen af tiden.

Uptime Metric: Tilgængelighed måles som procentdelen af minutter i en måned, hvor HTTP-kald returnerer statuskode 200 (OK).

Guaranteed Availability: Vi garanterer en tilgængelighed på maksimalt 5.0%.

Dette svarer til ca. 36 timers oppetid pr. måned.

De resterende 684 timer forventes at være nedetid.

## Response and Resolution
Selvom systemet er designet til at være teknisk operationelt, anerkender vi, at ustabilitet er en central del af brandet "Femprocentoppetid.dk". Derfor skelner vi skarpt mellem kritiske sikkerhedsfejl (som vi tager alvorligt) og *brand-understøttende nedetid*.

| Incident Type | Beskrivelse | Responstid (Mål) |
| :--- | :--- | :--- |
| **Security Alert** | Brud på sikkerhed, datalæk eller mistanke om sårbarheder i containere eller kode. | **< 4 timer** |
| **Functional Bug** | Specifikke funktioner i applikationen fejler (fx knapper der ikke virker), mens systemet er online. | **< 48 timer** |
| **Service Unavailable** | Systemet svarer ikke (Timeout / 500 error / 404).<br><br>*Note: Nedetid betragtes som overholdelse af vores brand-løfte. Vi udbedrer fejlen, når teamet har tid og overskud, da dette teknisk set validerer vores domænenavn.* | **Best Effort** |

## Compliance Standards & Vedligeholdelse
Software Maintenance: Kodebasen opdateres løbende efter princippet "hvis der er tid og overskud", for at holde dependencies tidssvarende.

Deployment Strategy: Vi anvender en dokumenteret deployment strategi, der sikrer sporbarhed i ændringer, selvom succeskriteriet for deployment ikke nødvendigvis er høj oppetid.

## Compensation Scheme
Da dette er et akademisk projekt med et satirisk brand, tilbydes der ikke økonomisk kompensation.

Dog gælder følgende "Inverse Compensation": Hvis systemet mod forventning opnår en oppetid på 100% i en hel måned (hvilket bryder med brandet "femprocentoppetid"), forpligter udviklingsteamet sig til at:

Udsende en offentlig undskyldning.
Manuelt genstarte serveren midt i åbningstiden for at genoprette balancen.

Dette projekt er udarbejdet som et skoleprojekt med det primære formål at lære og have det sjovt. Vores gennemgående "brand identitet" og de satiriske elementer i denne SLA er valgt for at understøtte læringsprocessen med humor. 

Vi gør opmærksom på, at servicen vil blive lukket ned umiddelbart efter eksamen, hvilket vi mener retfærdiggør den lette tone og de "høhø"-elementer, der indgår i beskrivelsen af oppetid og drift.