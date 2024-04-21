Feature: Create REST project
  In order to proceed fast working with microservices
  Me as a Golang developer
  Need to create a microservice project with kk cli commands
  Fast, easily and without copy & paste from existing projects 

  Background: Clean test environment
    Given a directory without folder godogs-rest-project
    And kk tool installed

  Scenario: Create chi based REST project
    When I create a chi project with uri github.com/waler4ik/godogs-rest-project
    Then godogs-rest-project contents are same as in cucumber/data/godogs-chi-project
  
  Scenario: Create gin based REST project
    When I create a gin project with uri github.com/waler4ik/godogs-rest-project
    Then godogs-rest-project contents are same as in cucumber/data/godogs-gin-project

  Scenario: Create chi REST project with websocket and endpoint 
    When I create a chi project with uri github.com/waler4ik/godogs-rest-project
    And I create a websocket with path /ws in folder godogs-rest-project
    And I create a resource with path /machines/data in folder godogs-rest-project
    And godogs-rest-project contents are same as in cucumber/data/godogs-chi-endpoints-project
    
  Scenario: Create gin REST project with websocket and endpoint 
    When I create a gin project with uri github.com/waler4ik/godogs-rest-project
    And I create a websocket with path /ws in folder godogs-rest-project
    And I create a resource with path /machines/data in folder godogs-rest-project
    And godogs-rest-project contents are same as in cucumber/data/godogs-gin-endpoints-project