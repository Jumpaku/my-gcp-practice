# Firebase Emulator Suite container
# Lightweight Node image
FROM node:20-alpine

# Install openjdk (for Firestore / Auth emulators) and curl
RUN apk add --no-cache openjdk17-jre-headless bash curl

# Set working directory
WORKDIR /usr/src/app

# Install Firebase CLI globally
RUN npm install -g firebase-tools@latest \
    && firebase --version

# Copy firebase.json and any emulator config from the project if present
# (these are optional; container can still run without them)
COPY firebase.json ./
COPY firestore.rules firestore.indexes.json storage.rules ./

# Expose default emulator ports (not all are required; expose the ones you use)
# UI: 4000
# Auth: 9099
# Firestore: 8080
# Realtime DB: 9000
# Functions: 5001
# Pub/Sub: 8085
# Storage: 9199
# Hosting: 5000
EXPOSE 4000 9099 8080 9000 5001 8085 9199 5000

# By default, run all configured emulators from firebase.json
# You can override this with `docker run ... firebase emulators:start --only firestore` etc.
CMD ["firebase", "emulators:start", "--project", "demo-project", "--only", "auth,firestore,functions,hosting,pubsub,storage", "--import", "/data", "--export-on-exit", "/data"]

