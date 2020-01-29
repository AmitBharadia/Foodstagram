# Foodstagram: Instagram for Food
### CMPE-281 Team Project - Team T800

## Team
- [Amit Bharadia](https://www.github.com/AmitBharadia)
- [Neil Shah](https://www.github.com/neilmshah)
- [Priyal Agrawal](https://github.com/Priyal-08)
- [Shabari Girish Ganapathy](https://www.github.com/shabari8695)

## Project Description
Foodstagram is a scalable containerized application similar to Instagram where users can post pictures of food for other users to view, like and comment. 

## Key Design Features
- CQRS
- Event Sourcing
- Graceful degradation of each microservice
- AKF scale cube for each microservice
  - X-axis - Load Balance between multiple cloned docker hosts
  - Y-axis - Microservice architecture to decouple functionalities
  - Z-axis - Database Sharding for each microservice

## System Architecture
![Foodstagram Architecture](https://github.com/nguyensjsu/fa19-281-t800/blob/master/screenshots/systemArchitecture.png)

## Foodstagram Demo
![Foodstagram Gif](https://github.com/nguyensjsu/fa19-281-t800/blob/master/screenshots/foodstagram.gif)
