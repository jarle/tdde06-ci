# Intro

*Visa upp flow.png*

*Förklara task management applikationen*

*Förklara jenkins och att jenkins + applikationen är två olika maskiner*

*Visa att all "kod" ligger på github*


# Demonstrering av applikationen

*hoppa in till maskinen*

*visa att allting är avstängt*

*Starta postgres container*

*Visa schema.sql*

*Lägg in schemat i databasen (genom docker)*

*Visa med `docker ps` att databasen är igång*

*Time to go! Kopiera in den långa raden och kör go applikationen!*

*kör igång nginx containern*

docker build -t nginxrunner -f nginx-docker.txt

*använd curl för att anropa applikationen och skapa lite listor och tasks*
curl -X POST -H "Content-Type: application/json" -d '{"Name": "Party"}' localhost/list
curl -X POST -H "Content-Type: application/json" -d '{"Name": "Balloons", "list_id": 2}' localhost/task
curl localhost/list
curl localhost/list/2

*run_everything.sh gör jobbet åt dig!*


# Demonstrering av git push

*Bryt builden och pusha!*

*Visa mail som säger att det ej fungerar*

*Fixa builden*

*Gå in i Jenkins och kolla att builden funkar*

*Visa Jenkins och tdde06 item och att det finns en nod*

*Visa build trigger i konfigurationen och visa webhook på github*

*Visa vad som körs (execute shell) och att ett mail kommer, om det går dåligt*
