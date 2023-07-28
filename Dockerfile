FROM node:alpine

# Update npm and clean cache
RUN apk update && apk upgrade && \
    apk add --no-cache npm && \
    rm -rf /var/cache/apk/*


# Install svgo and svgpath
RUN npm install -g svgo svgpath

# Export node modules path to NODE_PATH
ENV NODE_PATH /usr/local/lib/node_modules:$NODE_PATH

# Create shared directory
WORKDIR /shared
