# forum 

## Objectives

![Forum Screenshot](https://cdn.vox-cdn.com/thumbor/o52uRTrEg5vzSM-a-L92aESy_ds=/0x0:1409x785/1400x788/filters:focal(734x364:735x365)/cdn0.vox-cdn.com/uploads/chorus_asset/file/8846551/Screen_Shot_2017_07_13_at_1.09.20_PM.png)

This project is a web-based forum that facilitates communication between users, allowing them to create posts and comments, associate posts with categories, and like or dislike content. The forum also includes features like filtering posts by categories, created posts, and liked posts.

The backend is built using Go, and data is stored in an SQLite database. Docker is used to containerize the application for easy deployment.


## Building and Running the Application 

1. Clone the Repository
```
git clone https://01.alem.school/git/dabdrakhm/forum
```
2. Build the Docker Image:
```
make docker_build
```
3. Run the Docker Container:
```
make docker_run
```
4. Access the Forum: <br>
Open your web browser and navigate to http://localhost:4269 to start using the forum.

