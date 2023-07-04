
# Welcome To Go Lucky Wheel Game 
Lucky Wheel Game is Rest Api webservice based , it got daily limitation (3 times a day)
and it expires in 24 hour , and you can  play again in next day
## Software Architecture
This is an attempt to implement a clean architecture, and some other design patterns such as adapter and singelton in combination
## Requirements/Dependencies
- Docker
- Docker-compose
- golang:1.18-alpine docker image
- redis docker image
##  Getting Started
we have simple makefile in root of our project 
`make clean` 
will do everything for you to come up and runnig

## API Request
|      URL          |HTTP Method|Discription|
|----------------|-------------------------------|-----------------------------|
|`api/v1/lottery`|`POST`            |`Gives You a Prize `         |

body:
{
	UUID: “f1bc8f04-0500-11ee-be56-0242ac120002” // UUID of the user
}

or Whatever new user id requested
