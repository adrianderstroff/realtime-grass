# realtime-grass

![cover](/assets/images/github/cover.png)

This project is my attempt at Eddie Lee's [real-time grass](http://www.eddietree.com/#/grass/) demo, which was part of his master thesis. The demo consists of an infinite terrain covered in grass that waves in the wind that is exerted by the movement of the camera. Each grass blade is actual geometry as opposed to other approaches that use billboards to show grass patches. The goal is to have a decently fast simulation of natural looking grass.

To achieve this performance the load that is passed to the GPU on each frame has to be minimal. In addition level-of-detail (LOD) techniques are employed to adapt the detail of the grass blades depending on the distance to the camera.

## Skybox

Skybox **TropicalSunnyDay[Back|Down|Front|Left|Right|Up]** is taken from the project [skybox - a player skybox mod](http://minetest.daconcepts.com/my-main-mod-archive/sofars_mods/skybox/textures/) and renamed into **day/[back|bottom|front|left|right|top]**.

SkyboxSet by Heiko Irrgang ( http://gamvas.com ) is licensed under
the Creative Commons Attribution-ShareAlike 3.0 Unported License.
Based on a work at http://93i.de.


## Requirements
This project requires a GPU with OpenGL 4.3+ support.

The following dependencies depend on cgo. To make them work under Windows a compatible version of **mingw** is necessary. Information can be found [here](https://github.com/go-gl/glfw/issues/91). In my case I used *x86_64-7.2.0-posix-seh-rt_v5-rev1*. After installing the right version of **mingw** you can continue by installing the dependencies that follow next.

This project depends on **glfw** for creating a window and providing a rendering context, **go-gl/gl** for providing bindings to OpenGL and **go-gl/mathgl** provides vector and matrix math for OpenGL.
```
go get -u github.com/go-gl/glfw/v3.2/glfw
go get -u github.com/go-gl/gl/v4.3-core/gl
go get -u github.com/go-gl/mathgl/mgl32
```
After getting all dependencies the project should work without any errors.

## Theory

This section describes the idea behind the different parts of this project. It is to note that most of the theory is taken from Eddie Lee's masterthesis. Sections that are taken from somewhere else will be mentioned explicitly.

### Infinite terrain

As mentioned earlier does the demo contain a seemingly infinite terrain. However nothing truely infinite could be computed by the PC. Instead the landscape is going to be divided into smaller chunks that themselves contain a squared grid of tiles. Each tile consists of two triangles. To contour the landscape a height-map is used that is repeated infinitely.

![infinite terrain](/assets/images/github/infinite_terrain.png)

Now only a small portion of the terrain is shown at once. A radius *r_i* around the camera is used to determine how many chunks are loaded around the camera. Every frame all chunks that are outside the radius *r_i* are being destroyed. Next missing chunks that are now inside the radius *r_i* are being created.

When creating the chunk, a grid of tiles is created. For each tile the heights *h1, h2, h3, h4* of the four vertices that make up the tile are being taken from the height-map. For each tile the position of the tile and the plane data of both triangles that make of the tile are being stored. The plane equation is ![plane equation](/assets/images/github/plane.png) with **p** *= (x y z)* being a point on the plane, **n** *= (A B C)* the normal of the plane. 
![plane distance](/assets/images/github/plane-dist.png) 
is the distance of the plane from the origin. The normal of a plane can be calculated by taking the cross product between the vertices of the tile. The two normals of both planes are 
![plane normal 1](/assets/images/github/plane-normal1.png) 
for the upper right plane and 
![plane normal 2](/assets/images/github/plane-normal2.png)  
for the lower right plane.

To speed up the check for chunks that have to be created, the current chunk *(px, pz)* the camera is in, is calculated. Only x- and z-coordinate are relevant. Then the radius in number of chunks is calculated as 
![chunk radius](/assets/images/github/chunk-radius.png) with *t_c* being the side length of a chunk. Then iterating from *(px-cx, pz-cz)* to *(px+cx, pz+cz)* and taking the current *x* and *z* position as a hash for a map that maps strings onto chunks. If the x-z-coordinate is not in the map it must mean that the respective chunk does not exist yet. If it is not existent the distance of this chunk to the chunk where the camera resides in is checked and if the distance is smaller than *r_i* then this chunk gets created and the coordinate of this newly created chunk is added together with the chunk to the map.

### View Frustum Culling

The used approach is the clip-space View Frustum Culling approach from [lighthouse3d.com](http://www.lighthouse3d.com/tutorials/view-frustum-culling/clip-space-approach-extracting-the-planes/). Using the view and projection matrix of the camera AABBs can be transformed into clip-space. Checking whether an AABB is inside the View Frustum is as easy as checking if the points of the transformed AABB are inside the clip-space which is now a cube. To check if a point is inside the View Frustum the point has to be inside all planes defined by the six sides of the View Frustum cube. If the point is on the outside of one plane then the point is outside of the View Frustum. 
To speed things up we only check two sides of the AABB. The point that is pointing the most in direction of the plane's normal and the point that is on the opposite side of the first point. If one point of the AABB is inside the View Frustum and the other one is outside then the AABB gets intersected by the View Frustum.

Using View Frustum Culling only chunks that are either inside the View Frustum or intersect it are collected. The plane data of all selected chunks are uploaded into a vertex array buffer and are actually used. By this chunks that cannot be seen by the camera are not considered and thus precious computing time can be used on chunks that actually might appear inside the current frame.

### Terrain rendering

TODO

### Wind Simulation

TODO

### Grass rendering and simulation

TODO

### Postprocessing

TODO 

## ToDo

- [ ] Consistent variable and struct naming
- [ ] Using references where possible
- [ ] Rework Mesh code
- [ ] Rework Texture and Image
- [ ] Split Window and Interaction