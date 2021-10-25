/**
 * Parse the time elment's `datetime` attribute and set the local datetime as innerHTML.
 * @param {HTMLTimeElement} timeElement - the HTML time element to update
 */
function updateTimeElementToLocalTime(timeElement) {
  const datetimeAttribute = timeElement.getAttribute('datetime');
  if(typeof datetimeAttribute === 'string' && !!datetimeAttribute) {
    timeElement.innerHTML = new Date(datetimeAttribute).toLocaleString();
  }
}

export class GommentAdmin {
  /**
   * @constructor
   */
  constructor() {
    /** @type {HTMLTableElement | null} */
    this.table = null;
  }

  /**
   * Inject the gomment admin instance into a DOM.
   * @param {object} options - Object of DOM elements to inject into.
   * @param {string | HTMLTableElement} options.table - The table to display the comments in.
   */
  injectInto(options){
    if(typeof options.table === 'string') {
      const el = document.querySelector(options.table);
      if(!el || !(el instanceof HTMLTableElement)) {
        throw new Error(`HTML element with the specifier "${options.table}" was not found.`);
      }
      this.table = el;
    } else {
      this.table = options.table;
    }

    this.hydrate();
  }

  /**
   * Hydrate the server-rendered DOM with JS.
   */
  hydrate() {
    this.table.querySelectorAll('time').forEach(timeEl => updateTimeElementToLocalTime(timeEl));
  }
}
