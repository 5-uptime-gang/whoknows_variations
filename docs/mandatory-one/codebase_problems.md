# Problemer i kodebasen

## Kritiske problemer
- Koden kører på python 2 som stoppede med at få support d. 1. januar 2020. Skal konverteres til python 3 ✅
- Api key og database path ligge synligt i app.py
- No ID in the CREATE PAGES schema.sql
- SQL injectioner - den gamle kode bruger string interpolation til at køre sql statements. f.eks:
  
    `g.db.execute("SELECT id FROM users WHERE username = '%s'" % username)`
- 

## Vigtige problemer
- Change the footer, so it shows the current year dynamic

## Mindre problemer / kosmetiske 
- mordernized styling
- Styling doesnt work in mobile veiw
