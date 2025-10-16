# Generating Delivery Routes

There are thousand-to-millions coordinates that need to be processed into delivery routes. This is a Traveling Salesmans Problem (TSP), so generating the routes in realtime with this amount of data is very-very-**very** slow & resource expensive.

## Clustering to the Rescue

To solve this issue, we can generate multiple nested clusters and only choose the innermost cluster contains the addresses. Here are the visualization of the cluster (in a tree format, also imagine there are multiple country clusters).

![image](../asset/generate-route-cluster-example.svg)

> Recommendation:
>
> - Make each root cluster to have same level to make querying data easier.
> - Each innermost cluster have max 500-1.000 coordinate to optimize the time and resource needed to generate the route.

Then, system can generate routes for each cluster.

We need to generate the routes from country-to-country, city-by-city (inside a country), and next. With the resulting routes, Santa can start from the nearest country, then to the nearest city, then to the nearest district, and start deliver the present for each address in the district. After all presents in that district are sent, Santa will continue to the next district.

## Generating Distance

TSP algorithm requires a matrix containing the distance of each address, but we only have the coordinate of the addresses. To calculate the distance, we can use the Haversine distance formula.

## Actual Implementation

System must generate the route for each cluster and store it into a DB including the polyline data that needed to visualize the route. The client later can fetch needed routes.
