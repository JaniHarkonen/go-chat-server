import LoginView from "./view/loginView.js";

console.log("javascript successfully loaded");

var unmountPreviousView = () => {};

  // Clear the app root node and attach a new view via innerHTML
export function gotoView(view) {
  const hostElement = document.getElementById("app")

  if( !hostElement ) {
    console.log("ERROR: Unable to find host element to attach the view to!");
    return;
  }

    // Unmount the scripts of the previous view
  unmountPreviousView();

    // Mount the scripts of the new view
    // WARNING! THIS USE OF `window` IS DANGEROUS, HOWEVER, SUFFICIENT FOR THIS PROJECT
  for( let script of view.scripts ) {
    if( window[script.name] ) {
      throw new Erro("FATAL ERROR: Unable to mount script '" + script.name + "' due to overlap!");
    } else {
      window[script.name] = script;
    }
  }

  hostElement.innerHTML = view.html();
  view.onMount();

    // Create the next script unmounter
  unmountPreviousView = () => {
    for( let script of view.scripts ) {
      window[script.name] = undefined;
    }
  };
}

  // By default, goto login screen
gotoView(LoginView({id: "view-login"}));
