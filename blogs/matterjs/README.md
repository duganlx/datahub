# matter.js

文档: https://brm.io/matter-js/docs/

node 18

## 说明

Engine：The `Matter.Engine` module contains methods for creating and manipulating engines. An engine is a controller that manages updating the simulation of the world. See `Matter.Runner` for an optional game loop utility.

Render: The `Matter.Render` module is a simple canvas based renderer for visualising instances of `Matter.Engine`. It is intended for development and debugging purposes, but may also be suitable for simple games. It includes a number of drawing options including wireframe, vector with support for sprites and viewports.

Runner: The `Matter.Runner` module is an optional utility which provides a game loop, that handles continuously updating a `Matter.Engine` for you within a browser. It is intended for development and debugging purposes, but may also be suitable for simple games. If you are using your own game loop instead, then you do not need the `Matter.Runner` module. Instead just call `Engine.update(engine, delta)` in your own loop.

Body: The `Matter.Body` module contains methods for creating and manipulating rigid bodies. For creating bodies with common configurations such as rectangles, circles and other polygons see the module `Matter.Bodies`.

Bodies: The `Matter.Bodies` module contains factory methods for creating rigid body models with commonly used body configurations (such as rectangles, circles and other polygons).

Composite: A composite is a collection of `Matter.Body`, `Matter.Constraint` and other `Matter.Composite` objects. They are a container that can represent complex objects made of multiple parts, even if they are not physically connected. A composite could contain anything from a single body all the way up to a whole world. When making any changes to composites, use the included functions rather than changing their properties directly.

Composites: The `Matter.Composites` module contains factory methods for creating composite bodies with commonly used configuration (such as stacks and chains).

Constraint: The `Matter.Constraint` module contains methods for creating and manipulating constraints. Constraints are used for specifying that a fixed distance must be maintained between two bodies (or a body and a fixed world-space position). The stiffness of constraints can be modified to create springs or elastic.

MouseConstraint: The `Matter.MouseConstraint` module contains methods for creating mouse constraints. Mouse constraints are used for allowing user interaction, providing the ability to move bodies via the mouse of touch.

Events: The `Matter.Events` module contains methods to fire and listen to events on other objects.

Common: The `Matter.Common` module contains utility functions that are common to all modules.

Plugin: The `Matter.Plugin` module contains functions for registering and installing plugins on modules.
