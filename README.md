# Vencom

Designed using [clonemap](https://git.rwth-aachen.de/acs/public/cloud/mas/clonemap), a project of RWTH AACHEN University

## Installation

Deploy clonemap's core docker containers:

```cmd
git clone https://git.rwth-aachen.de/acs/public/cloud/mas/clonemap.git "<folder>"
"<folder>"\deployments\docker docker-compose up -d #REM to start the containers
"<folder>"\deployments\docker docker-compose down
```

Build Vencom image and deploy docker container/s:

```cmd
git clone https://github.com/SSII-UEM/Vencom.git "<folder>"
cd "<folder>"
docker build --rm -f build\Dockerfile -t vencom .
curl -X "POST" -d @scenario.json localhost:30009/api/clonemap/mas
```

## Issues

* Vencom is currently designed to work on a single MAS, so configuring more `imagegroups.agents` than the `config.agentsperagency` will not work as expected.

* Buyers bought products aren't always properly recovered

* Buyers HP can reach subzero values, agents never die

* Periodic behaviours introduce racing behaviours while reading/writing (concurrent) each agent properties

* Coordinate distances aren't properly computed

## Future additions

* Products should have a lifespan to incentivate retailers to design better market strategies and buyers to eat them before others they own

* Buyers should be able to generate an income by themselves in order to survive. That means working, stealing, gifting...

* Agents are located on a 3D map, so it would be interesting to make them able to move (retailers to find better spots to sell, buyers better spots in order to perform their activities)

* Retailers should also have the same capabilities as buyers. That means they should also need to consume products to survive.

* Intelligence