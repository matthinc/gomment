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

class EventBus {
  constructor() {
    this.subscribers = [];
  }

  on(eventName, callback) {
    this.subscribers.push({eventName, callback});
  }

  emit(eventName, eventData) {
    this.subscribers
      .filter(s => s.eventName === eventName)
      .forEach(s => s.callback(eventData));
  }
}

class ApiClient {
  /**
   * Delete the comment with the given id
   * @param {string} commentId - comment id to delete
   * @returns {Promise<void>}
   */
  async deleteComment(commentId) {
    return new Promise((resolve, reject) => {
      console.log(`stub: delete comment ${commentId}`);
      resolve();
    });
  }
}

class CommentRow {
  /**
   * @constructor
   * @param {any} tableRow - a HTML table row containing one comment
   * @param {EventBus} bus - simple event bus
   */
  constructor(tableRow, bus) {
    if(!(tableRow instanceof HTMLTableRowElement)) {
      throw new Error('the provided element is not a valid HTMLTableRowElement');
    }
    /** @type HTMLTableRowElement */
    this.tableRow = tableRow;

    /** @type EventBus */
    this.bus = bus;

    const id = tableRow.getAttribute('data-id');
    if(typeof id !== 'string') {
      throw new Error('CommentRow does not contain valid id');
    }
    /** @type string */
    this.id = id;

    /** @type boolean */
    this.isEdit = false;

    this.hydrate();
  }

  hydrate() {
    this.tableRow.querySelector('a[href="#delete"]').addEventListener('click', e => {
      e.preventDefault();
      this.bus.emit('delete', this.id);
    });

    this.tableRow.querySelector('a[href="#edit"]').addEventListener('click', e => {
      e.preventDefault();
      this.enableEditMode();
    });
  }

  /**
   * Get the id of the comment belonging to this row.
   * @returns {string}
   */
  getId() {
    return this.id;
  }

  /**
   * Enable the edit mode
   */
  enableEditMode() {
    if(this.isEdit) {
      return;
    }

    /** @type HTMLParagraphElement */
    const commentParagraph = this.tableRow.querySelector('.gmt__comment p');
    const parent = commentParagraph.parentElement;
    const commentText = commentParagraph.textContent;
    const commentTextarea = document.createElement('textarea');
    commentTextarea.innerText = commentText;
    parent.removeChild(commentParagraph);
    parent.appendChild(commentTextarea);

    this.isEdit = true;
  }
}

export class GommentAdmin {
  /**
   * @constructor
   */
  constructor() {
    /** @type {HTMLTableElement | null} */
    this.table = null;

    /** @type {ApiClient} */
    this.api = new ApiClient();

    /** @type {EventBus} */
    this.bus = new EventBus();
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

    this.bus.on('delete', (commentId) => {
      this.deleteComment(commentId);
    });

    this.hydrate();
  }

  /**
   * Hydrate the server-rendered DOM with JS.
   */
  hydrate() {
    this.table.querySelectorAll('time').forEach(timeEl => updateTimeElementToLocalTime(timeEl));
    this.table.querySelectorAll('tr.gmt__row').forEach(row => {
      new CommentRow(row, this.bus);
    });
  }

  /**
   * Delete a comment, ask for confirmation before
   * @param {string} commentId
   */
  deleteComment(commentId) {
    this.api.deleteComment(commentId)
      .then(() => {
        console.log('delete ok');
      })
      .catch((err) => {
        console.error(err);
      });
  }
}
