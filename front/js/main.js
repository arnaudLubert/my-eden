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
arObject.visible = false;
scene.addNode(arObject);

let flower = new Gltf2Node({url: '/models/gltf/random/cylinder.glb'});
flower.scale = [0.02, 0.02, 0.1]; // custom scale
flower.rotation[0] = -0.7;

arObject.addNode(flower);

let reticle = new Gltf2Node({url: '/models/gltf/reticle/reticle.gltf'});
reticle.visible = false;
scene.addNode(reticle);

// Having a really simple drop shadow underneath an object helps ground
// it in the world without adding much complexity.
let shadow = new DropShadowNode();
vec3.set(shadow.scale, 0.15, 0.15, 0.15);
arObject.addNode(shadow);

const MAX_FLOWERS = 4;
let flowers = [];

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

    if (flowers.length >= MAX_FLOWERS) {
     // let oldFlower = flowers.shift();
     scene.removeNode(flowers[flowerCounter]);
     flowers[flowerCounter] = newFlower;
      flowerCounter++;

      if (flowerCounter === MAX_FLOWERS)
        flowerCounter = 0;
  } else
    flowers.push(newFlower);
  globalMat = newFlower;

  if (flowers.length > 1) {
     // const distance = Math.sqrt(Math.pow(flowers[flowers.length - 2].matrix[12] - flowers[flowers.length - 1].matrix[12], 2) + Math.pow(flowers[flowers.length - 2].matrix[14] - flowers[flowers.length - 1].matrix[14], 2));
    //  console.log(distance);
  }
 //console.log(globalMat.matrix);

 if (flowers.length == MAX_FLOWERS) {
     const corners = [0, 1, 2, 3];
     const ab = Math.sqrt(Math.pow(flowers[corners[0]].matrix[12] - flowers[corners[1]].matrix[12], 2) + Math.pow(flowers[corners[0]].matrix[14] - flowers[corners[1]].matrix[14], 2));
     const bc = Math.sqrt(Math.pow(flowers[corners[1]].matrix[12] - flowers[corners[2]].matrix[12], 2) + Math.pow(flowers[corners[1]].matrix[14] - flowers[corners[2]].matrix[14], 2));
     const cd = Math.sqrt(Math.pow(flowers[corners[2]].matrix[12] - flowers[corners[3]].matrix[12], 2) + Math.pow(flowers[corners[2]].matrix[14] - flowers[corners[3]].matrix[14], 2));
     const ad = Math.sqrt(Math.pow(flowers[corners[3]].matrix[12] - flowers[corners[0]].matrix[12], 2) + Math.pow(flowers[corners[3]].matrix[14] - flowers[corners[0]].matrix[14], 2));

     const angleA = find_angle(
         {x: flowers[corners[1]].matrix[12], y: flowers[corners[1]].matrix[14]},
         {x: flowers[corners[3]].matrix[12], y: flowers[corners[3]].matrix[14]},
         {x: flowers[corners[0]].matrix[12], y: flowers[corners[0]].matrix[14]});
     const angleC = find_angle(
         {x: flowers[corners[3]].matrix[12], y: flowers[corners[3]].matrix[14]},
         {x: flowers[corners[1]].matrix[12], y: flowers[corners[1]].matrix[14]},
         {x: flowers[corners[2]].matrix[12], y: flowers[corners[2]].matrix[14]});

//    console.log("angles:",angleA * 180 / Math.PI, angleC * 180 / Math.PI);
    const surface = 0.5 * ab * ad * Math.sin(angleA) + 0.5 * bc * cd * Math.sin(angleC); // m2
    console.log("surface:", Math.round(surface * 100) / 100, 'm2');
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
