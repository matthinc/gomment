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

    this.render();
  }

  /**
   * Render the table based on the current filters and settings.
   */
  render() {

  }
}
