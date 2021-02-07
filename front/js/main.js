import {WebXRButton} from '/js/util/webxr-button.js';
import {Scene} from '/js/render/scenes/scene.js';
import {Renderer, createWebGLContext} from '/js/render/core/renderer.js';
import {Node} from '/js/render/core/node.js';
import {Gltf2Node} from '/js/render/nodes/gltf2.js';
import {DropShadowNode} from '/js/render/nodes/drop-shadow.js';
import {SevenSegmentText} from '/js/render/nodes/seven-segment-text.js';
import {vec3} from '/js/render/math/gl-matrix.js';
import {Ray} from '/js/render/math/ray.js';

// XR globals.
let xrButton = null;
let xrRefSpace = null;
let xrViewerSpace = null;
let xrHitTestSource = null;

// WebGL scene globals.
let gl = null;
let renderer = null;
let scene = new Scene();
let flowerCounter = 0;
scene.enableStats(false);

let arObject = new Node();
let flowerPrefab = new Node();
arObject.visible = false;
flowerPrefab.visible = false;
scene.addNode(arObject);
scene.addNode(flowerPrefab);

let flower = new Gltf2Node({url: '/models/gltf/sunflower/sunflower.gltf'});
let stick = new Gltf2Node({url: '/models/gltf/random/cylinder.glb'});
stick.scale = [0.02, 0.02, 0.1]; // custom scale
stick.rotation[0] = -0.7;


arObject.addNode(stick);
flowerPrefab.addNode(flower);

let reticle = new Gltf2Node({url: '/models/gltf/reticle/reticle.gltf'});
reticle.visible = false;
scene.addNode(reticle);

// Having a really simple drop shadow underneath an object helps ground
// it in the world without adding much complexity.
let shadow = new DropShadowNode();
vec3.set(shadow.scale, 0.15, 0.15, 0.15);
arObject.addNode(shadow);
flowerPrefab.addNode(shadow);

const MAX_CORNERS = 4;
let corners = [];
let garden = [];

// Ensure the background is transparent for AR.
scene.clear = false;

function initXR() {
  xrButton = new WebXRButton({
    onRequestSession: onRequestSession,
    onEndSession: onEndSession,
    textEnterXRTitle: "CREATE MY EDEN",
    textXRNotFoundTitle: "AR NOT FOUND",
    textExitXRTitle: "EXIT  AR",
  });
  document.querySelector('header').appendChild(xrButton.domElement);

  if (navigator.xr) {
    navigator.xr.isSessionSupported('immersive-ar')
                .then((supported) => {
      xrButton.enabled = supported;
    });
  }
}

function onRequestSession() {
  return navigator.xr.requestSession('immersive-ar', {requiredFeatures: ['local', 'hit-test']})
                     .then((session) => {
    xrButton.setSession(session);
    onSessionStarted(session);
  });
}

function onSessionStarted(session) {
  session.addEventListener('end', onSessionEnded);
  session.addEventListener('select', onSelect);

  if (!gl) {
    gl = createWebGLContext({
      xrCompatible: true
    });

    document.body.appendChild(gl.canvas);

    function onResize() {
      gl.canvas.width = gl.canvas.clientWidth * window.devicePixelRatio;
      gl.canvas.height = gl.canvas.clientHeight * window.devicePixelRatio;
    }
    window.addEventListener('resize', onResize);
    onResize();

    renderer = new Renderer(gl);

// GUI
/*
    let uiText = new SevenSegmentText();
    // Hard coded because it doesn't change:
    // Scale by 0.075 in X and Y
    // Translate into upper left corner w/ z = 0.02
    uiText.matrix = new Float32Array([
        0.075, 0, 0, 0,
        0, 0.075, 0, 0,
        0, 0, 1, 0,
        -0.3625, 0.3625, 0.02, 1,
    ]);
    uiText.onRendererChanged(renderer);
    uiText.text = "1234556FPS";
    scene.addNode(uiText);
*/
    scene.setRenderer(renderer);

  }

  session.updateRenderState({ baseLayer: new XRWebGLLayer(session, gl) });

  // In this sample we want to cast a ray straight out from the viewer's
  // position and render a reticle where it intersects with a real world
  // surface. To do this we first get the viewer space, then create a
  // hitTestSource that tracks it.
  session.requestReferenceSpace('viewer').then((refSpace) => {
    xrViewerSpace = refSpace;
    session.requestHitTestSource({ space: xrViewerSpace }).then((hitTestSource) => {
      xrHitTestSource = hitTestSource;
    });
  });

  session.requestReferenceSpace('local').then((refSpace) => {
    xrRefSpace = refSpace;

    session.requestAnimationFrame(onXRFrame);
  });
}

function onEndSession(session) {
  xrHitTestSource.cancel();
  xrHitTestSource = null;
  session.end();
}

function onSessionEnded(event) {
  xrButton.setSession(null);
}

// Adds a new object to the scene at the
// specificed transform.
function addARObjectAt(matrix) {
    let newFlower = arObject.clone();
    newFlower.visible = true;
    newFlower.matrix = matrix;
    scene.addNode(newFlower);

    if (corners.length >= MAX_CORNERS) {
     // let oldFlower = corners.shift();
     scene.removeNode(corners[flowerCounter]);
     corners[flowerCounter] = newFlower;
      flowerCounter++;

      if (flowerCounter === MAX_CORNERS)
        flowerCounter = 0;
  } else
    corners.push(newFlower);
//  globalMat = newFlower;
/*
  if (corners.length > 1) {
      const distance = Math.sqrt(Math.pow(corners[selector.length - 2].matrix[12] - corners[selector.length - 1].matrix[12], 2) + Math.pow(corners[selector.length - 2].matrix[14] - corners[selector.length - 1].matrix[14], 2));
      console.log(distance);
}*/
 //console.log(globalMat.matrix);

 if (corners.length == MAX_CORNERS) {
     const selector = [0, 1, 2, 3];
     const ab = Math.sqrt(Math.pow(corners[selector[0]].matrix[12] - corners[selector[1]].matrix[12], 2) + Math.pow(corners[selector[0]].matrix[14] - corners[selector[1]].matrix[14], 2));
     const bc = Math.sqrt(Math.pow(corners[selector[1]].matrix[12] - corners[selector[2]].matrix[12], 2) + Math.pow(corners[selector[1]].matrix[14] - corners[selector[2]].matrix[14], 2));
     const cd = Math.sqrt(Math.pow(corners[selector[2]].matrix[12] - corners[selector[3]].matrix[12], 2) + Math.pow(corners[selector[2]].matrix[14] - corners[selector[3]].matrix[14], 2));
     const ad = Math.sqrt(Math.pow(corners[selector[3]].matrix[12] - corners[selector[0]].matrix[12], 2) + Math.pow(corners[selector[3]].matrix[14] - corners[selector[0]].matrix[14], 2));

     const angleA = find_angle(
         {x: corners[selector[1]].matrix[12], y: corners[selector[1]].matrix[14]},
         {x: corners[selector[3]].matrix[12], y: corners[selector[3]].matrix[14]},
         {x: corners[selector[0]].matrix[12], y: corners[selector[0]].matrix[14]});
     const angleC = find_angle(
         {x: corners[selector[3]].matrix[12], y: corners[selector[3]].matrix[14]},
         {x: corners[selector[1]].matrix[12], y: corners[selector[1]].matrix[14]},
         {x: corners[selector[2]].matrix[12], y: corners[selector[2]].matrix[14]});

//    console.log("angles:",angleA * 180 / Math.PI, angleC * 180 / Math.PI);
    const surface = 0.5 * ab * ad * Math.sin(angleA) + 0.5 * bc * cd * Math.sin(angleC); // m2
//    console.log("surface:", Math.round(surface * 100) / 100, 'm2');
    document.getElementById('surface').children[0].textContent = 'Surface: ' + (Math.round(surface * 100) / 100) + ' m2';

    matrix[14] = corners[selector[0]].matrix[14];
    console.log(corners[selector[2]].matrix[14], matrix[14]);

    while (matrix[14] < corners[selector[2]].matrix[14]) {
        matrix[12] = corners[selector[0]].matrix[12];
        while (matrix[12] < corners[selector[1]].matrix[12]) {
            newFlower = flowerPrefab.clone();
            newFlower.visible = true;
            newFlower.matrix = matrix;
            garden.push(newFlower);
            scene.addNode(newFlower);
            matrix[12] += 0.3;
        }
        matrix[14] += 0.3;
    }
//corners[selector[0]].matrix[12] + 0.2;
//corners[selector[0]].matrix[14] + 0.2;
 }
}

function find_angle(p0,p1,c) {
    var p0c = Math.sqrt(Math.pow(c.x-p0.x,2)+
                        Math.pow(c.y-p0.y,2)); // p0->c (b)
    var p1c = Math.sqrt(Math.pow(c.x-p1.x,2)+
                        Math.pow(c.y-p1.y,2)); // p1->c (a)
    var p0p1 = Math.sqrt(Math.pow(p1.x-p0.x,2)+
                         Math.pow(p1.y-p0.y,2)); // p0->p1 (c)
    return Math.acos((p1c*p1c+p0c*p0c-p0p1*p0p1)/(2*p1c*p0c));
}

let rayOrigin = vec3.create();
let rayDirection = vec3.create();

function onSelect(event) {
  if (reticle.visible) {
    // The reticle should already be positioned at the latest hit point,
    // so we can just use it's matrix to save an unnecessary call to
    // event.frame.getHitTestResults.
    addARObjectAt(reticle.matrix);
  }
}

// Called every time a XRSession requests that a new frame be drawn.
function onXRFrame(t, frame) {
  let session = frame.session;
  let pose = frame.getViewerPose(xrRefSpace);

  reticle.visible = false;

  // If we have a hit test source, get its results for the frame
  // and use the pose to display a reticle in the scene.
  if (xrHitTestSource && pose) {
    let hitTestResults = frame.getHitTestResults(xrHitTestSource);
    if (hitTestResults.length > 0) {
      let pose = hitTestResults[0].getPose(xrRefSpace);
      reticle.visible = true;
      reticle.matrix = pose.transform.matrix;
    }
  }

  scene.startFrame();

  session.requestAnimationFrame(onXRFrame);

  scene.drawXRFrame(frame, pose);

  scene.endFrame();
}

// Start the XR application.
initXR();
/*
12=X/Est
13=Z
14=Y/Nord

178.33.193.240

*/
