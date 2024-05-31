# Utilisez une image de base contenant Go pour compiler et exécuter le projet
FROM golang:latest

# Définissez le répertoire de travail à l'intérieur du conteneur
WORKDIR /app

# Copiez les fichiers de module Go et téléchargez les dépendances
COPY go.sum go.sum
COPY go.mod go.mod
RUN go mod download

# Copiez le reste du code source
COPY main.go main.go
COPY ./bet ./bet
COPY ./logger ./logger

# Expose le port sur lequel l'application écoute
EXPOSE 9000

# Commande à exécuter lors du démarrage du conteneur
CMD ["go", "run", "main.go", "9000"]
