# This is a basic workflow to help you get started with Actions

name: Push on tags

# Controls when the workflow will run
on:
  create:
    tags:
      - "v*"

# A workflow run is made up of one or more jobs that can run sequentially or in parallel
jobs:
  # This workflow contains a single job called "build"
  build:
    # The type of runner that the job will run on
    runs-on: ubuntu-latest

    # Steps represent a sequence of tasks that will be executed as part of the job
    steps:
      # Checks-out your repository under $GITHUB_WORKSPACE, so your job can access it
      - uses: actions/checkout@v2
      
      # Runs a set of commands using the runners shell
      - name: Build the Docker image
        run: |
          make build
          docker images

      # Login with docker
      - name: Login to ghcr.io
        run: |
          echo ${{ secrets.GITHUB_TOKEN }} | docker login ghcr.io -u seeker89 --password-stdin
          make tag
          make push
