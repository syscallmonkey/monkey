name: Push on tags
on:
  create:
    tags:
      - "v*"
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Build the Docker image
        run: |
          make build
          docker images
      - name: Login to ghcr.io
        run: |
          echo ${{ secrets.GITHUB_TOKEN }} | docker login ghcr.io -u seeker89 --password-stdin
          make tag
          make push
