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

  /**
   * Edit the comment with the given id by replacing the text
   * @param {string} commentId - comment id to edit
   * @param {string} text - new comment text
   * @returns {Promise<void>}
   */
  async editComment(commentId, text) {
    return new Promise((resolve, reject) => {
      console.log(`stub: edit comment ${commentId}: ${text}`);
      setTimeout(()=>resolve(), 1000);
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
    this.setupActions('default');
  }

  /**
   * Setup the action buttons according to a given state
   * @param {'none' | 'default' | 'delete'} state - The state for which to setup the buttons
   */
  setupActions(state) {
    const actions = this.tableRow.querySelector('.gmt__actions');

    switch(state) {
    case 'none':
      actions.innerHTML = '';
      break;
    case 'default':
      actions.innerHTML = '<a href="#edit">Edit</a> | <a href="#delete">Delete</a>';

      this.tableRow.querySelector('a[href="#delete"]').addEventListener('click', e => {
        e.preventDefault();
        this.bus.emit('delete', this);
      });

      this.tableRow.querySelector('a[href="#edit"]').addEventListener('click', e => {
        e.preventDefault();
        this.enableEditMode();
      });
      break;
    case 'delete':
      actions.innerHTML = '<a href="#confirm">Confirm</a> | <a href="#abort">Abort</a>';

      this.tableRow.querySelector('a[href="#abort"]').addEventListener('click', e => {
        e.preventDefault();
        this.disableEditMode(true);
      });

      this.tableRow.querySelector('a[href="#confirm"]').addEventListener('click', e => {
        e.preventDefault();
        this.disableEditMode(false);
        this.setupActions('none');
        this.setLoading(true);
        this.bus.emit('edit', this);
      });
      break;
    }
  }

  /**
   * Get the id of the comment belonging to this row.
   * @returns {string}
   */
  getId() {
    return this.id;
  }

  /**
   * Get the current text contents (not HTML escaped)
   * @returns {string}
   */
  getText() {
    if(this.isEdit) {
      return this.tableRow.querySelector('.gmt__comment textarea').value;
    } else {
      return this.tableRow.querySelector('.gmt__comment p').innerText;
    }
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
    this.commentText = commentParagraph.innerText;
    const commentTextarea = document.createElement('textarea');
    commentTextarea.innerText = this.commentText;
    parent.removeChild(commentParagraph);
    parent.appendChild(commentTextarea);

    this.setupActions('delete');

    this.isEdit = true;
  }

  /**
   * Disable the edit mode
   * @param {boolean} discard - true if the modifications shall be discarded
   */
  disableEditMode(discard) {
    if(!this.isEdit) {
      return;
    }

    /** @type HTMLTextAreaElement */
    const commentTextarea = this.tableRow.querySelector('.gmt__comment textarea');
    const parent = commentTextarea.parentElement;
    let commentText;
    if(discard) {
      commentText = this.commentText;
      delete this.commentText;
    } else {
      commentText = commentTextarea.value;
    }
    const commentParagraph = document.createElement('p');
    commentParagraph.innerText = commentText;
    parent.removeChild(commentTextarea);
    parent.appendChild(commentParagraph);

    this.isEdit = false;
  }

  /**
   * Visually indicate some async action is being done on the row.
   * @param {boolean} loading - true if the loading state shall be set
   */
  setLoading(loading) {
    if(loading) {
      this.tableRow.classList.add('loading');
    } else {
      this.tableRow.classList.remove('loading');
    }
  }

  removeSelf() {
    this.tableRow.parentElement.removeChild(this.tableRow);
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

    this.bus.on('delete', this.deleteComment.bind(this));
    this.bus.on('edit', this.editComment.bind(this));

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
   * @param {CommentRow} commentRow
   */
  deleteComment(commentRow) {
    this.api.deleteComment(commentRow.getId())
      .then(() => {
        console.log('delete ok');
        commentRow.removeSelf();
      })
      .catch((err) => {
        console.error(err);
      });
  }

  /**
   * Edit a comment, ask for confirmation before
   * @param {CommentRow} commentRow
   */
  editComment(commentRow) {
    this.api.editComment(commentRow.getId(), commentRow.getText())
      .then(() => {
        console.log('edit ok');
        commentRow.setLoading(false);
        commentRow.setupActions('default');
      })
      .catch((err) => {
        console.error(err);
      });
  }
}
