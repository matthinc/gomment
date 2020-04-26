/**
 * A response returned by the backend for a comment query.
 * @typedef {Object} CommentQueryResponse
 * @property {Array<Comment>} comments The queried comments.
 * @property {number} total The amount of comments returned by the query.
 */

/**
 * Content of a single comment.
 * @typedef {Object} CommentData
 * @property {string} author The author of the comment.
 * @property {number} created_at Creation date of the comment.
 * @property {string} text The content of the comment.
 */

/**
 * A comment response from the backend.
 * @typedef {Object} Comment
 * @property {CommentData} [comment] The actual comment data (content).
 * @property {Array<Comment>} [children] Children comments of this comment.
 * @property {boolean} [has_children] Whether the comment has (unloaded) child comments.
 */

/**
 * Creates a new element and inserts it as the child of parent.
 * @param {string} type - Type of the new DOM element.
 * @param {string} className - Name of the class to assign to the newly created element.
 * @param {HTMLElement} parent - The new parent of the DOM element.
 * @param {*} attributes - Attributes to assign to the new DOM element.
 * @returns {HTMLElement} The newly created child.
 */
function insertElement(type, className, parent, attributes = {}) {
  if (parent === undefined) {
    throw new Error('No valid parent element was provided.');
  }
  
  const elem = document.createElement(type);
  elem.className = className;
  // Misc. attributes
  for (const key in attributes) {
    // @ts-ignore
    elem[key] = attributes[key];
  }
  parent.appendChild(elem);
  return elem;
}

export class Gomment {

  /**
   * @constructor
   * @param {?*} [options]
   */
  constructor(options) {
    if (!options || typeof options !== 'object') {
      throw new Error('mandatory parameter `options` was not provided');
    }
    
    // Required options
    /** @type {string} */
    this.apiURL = options.apiURL;
    // append a missing trailing backslash
    this.apiURL = this.apiURL.endsWith('/') ? this.apiURL : `${this.apiURL}/`;
    
    /** @type {string} */
    this.thread = options.thread;

    // Optional options
    /** @type {number} */
    this.batchSize = options.batchSize || 10;
    /** @type {number} */
    this.maxDepth = options.maxDepth || 2;

    // i18n
    this.i18n = options.i18n || {
      title: 'Comments',
      empty: 'No comments',
      input_title: 'Leave a comment:',
      placeholder_name: 'Name',
      placeholder_email: 'E-Mail',
      placeholder_text: 'Your comment',
      submit: 'Submit',
      show_more: 'Load more comments',
      show_more_depth: 'Load more replies',
      /** @type {(date: Date) => string} */
      format_date: (date) => `${date.getFullYear()}.${date.getMonth()}.${date.getDate()} ${date.getHours()}:${date.getMinutes()}`
    };

    // stateful information
    /** @type {?} */
    this.comments = [];
    /** @type {number} */
    this.lastOffset = 0;
    /** @type {HTMLElement | null} */
    this.commentsElement = null;
  }

  /**
   * Query for comments with the specified parameters.
   * @param {string} threadId - Id of the comment thread.
   * @param {number} offset - The offset.
   * @param {number} max - Maximum amount of comments to query.
   * @param {number} depth - Maximum level of depth to query.
   * @returns {Promise<Response>} - HTML Response.
   */
  queryComments(threadId, offset, max, depth) {
    return window.fetch(`${this.apiURL}comments?thread=${threadId}&offset=${offset}&max=${max}&depth=${depth}`);
  }

  /**
   * Query for comments with the specified parameters and render them.
   * @param {number} offset - The offset.
   * @param {number} max - Maximum amount of comments to query.
   * @param {number} depth - Maximum level of depth to query.
   * @returns {void}
   */
  loadComments(offset, max, depth) {
    /** @type {(data: CommentQueryResponse) => void} */
    const handler = (data) => {
      data.comments.forEach(
        /** @type {(item: Comment, index: number) => any} */
        (item, index) => this.comments[index + offset] = item
      );
      this.lastOffset = offset;
      this.renderComments(this.comments, data.total);
    };

    this.queryComments(this.thread, offset, max, depth)
      .then(data => data.json())
      .then(handler);
  }

  /**
   * Load comments with the initial default parameters.
   * @returns {void}
   */
  loadCommentsInitial() {
    this.loadComments(0, this.batchSize, this.maxDepth);
  }

  /**
   * Load 'more' comments.
   * @returns {void}
   */
  loadNextBatch() {
    this.loadComments(this.lastOffset + this.batchSize, this.batchSize, this.maxDepth);
  }

  /**
   * Load 'more' replies.
   * @param {number} index - The comment to load more replies from.
   * @param {number} depth - Maximum level of depth to query.
   * @returns {void}
   */
  loadMoreDepth(index, depth) {
    this.loadComments(index, 1, depth + this.batchSize);
  }

  /**
   * Render single comment and children recursively.
   * @param {HTMLElement} parent - The comment tree parent.
   * @param {Comment} comment - The offset.
   * @param {number} treeIndex - The relative offset in the tree.
   * @param {number} depth - Maximum level of depth to query.
   * @returns {void}
   */
  renderComment(parent, comment, treeIndex, depth) {
    const commentElement = insertElement('div', 'gomment-comment', parent);
    insertElement('div', 'gomment-comment-author', commentElement, { innerHTML: comment.comment.author });
    insertElement('div', 'gomment-comment-date', commentElement, { innerHTML: this.i18n.format_date(new Date(comment.comment.created_at)) });
    insertElement('div', 'gomment-comment-text', commentElement, { innerHTML: comment.comment.text });
    const childrenElement = insertElement('div', 'gomment-comment-children', commentElement);

    if (comment.children) {
      for (const childComment of comment.children) {
        this.renderComment(childrenElement, childComment, treeIndex, depth + 1);
      }
    } else if (comment.has_children) {
      // No children but hasChildren -> Load more button
      insertElement('button', 'gomment-show-more-depth-button', childrenElement, {
        innerHTML: this.i18n.show_more_depth,
        onclick: () => this.loadMoreDepth(treeIndex, depth)
      });
    }
  }

  /**
   * Render all (top level) comments with children.
   * @param {Array<Comment>} comments - The comments to render.
   * @param {number} total
   * @returns {void}
   */
  renderComments(comments, total) {
    // Remove all comments from element
    if (this.commentsElement === null) {
      console.warn('Gomment instance needs to be injected before rendering.');
      return;
    }
    /** @type {HTMLElement} */
    const ce = this.commentsElement;
    
    ce.innerHTML = '';

    // Render comments
    comments.forEach((comment, index) => this.renderComment(ce, comment, index, 0));
    // Show more button
    if (comments.length < total) {
      const showMoreContainer = insertElement('div', 'gomment-show-more-container', ce);
      insertElement('button', 'gomment-show-more-button', showMoreContainer, {
        innerHTML: this.i18n.show_more,
        onclick: () => this.loadNextBatch()
      });
    }
  }

  /**
   * Injects the gomment instance into a parent element.
   * @param {string | HTMLElement} element - Parent element.
   */
  injectInto(element) {
    /** @type {HTMLElement} */
    let container;
    if (typeof element === 'string') {
      const el = document.querySelector(element);
      if (!el || !(el instanceof HTMLElement)) {
        throw new Error('HTML element with the specifier "${element}" was not found.');
      }
      container = insertElement('div', 'gomment', el);
    } else if (element instanceof HTMLElement) {
      container = element;
    } else {
      throw new Error('Parent element needs to be provided as a query selector or `HTMLElement`');
    }
    
    // Create container element
    insertElement('div', 'gomment-title', container, { innerHTML: this.i18n.title });

    // Comment input section
    /** @type {HTMLElement} */
    const inputSection = insertElement('div', 'gomment-input-section', container);
    insertElement('div', 'gomment-input-title', inputSection, { innerHTML: this.i18n.input_title });
    insertElement('input', 'gomment-email', inputSection, { placeholder: this.i18n.placeholder_email });
    insertElement('input', 'gomment-display-name', inputSection, { placeholder: this.i18n.placeholder_name });
    insertElement('textarea', 'gomment-text-input', inputSection, { placeholder: this.i18n.placeholder_text });
    insertElement('button', 'gomment-submit-button', inputSection, { innerHTML: this.i18n.submit });

    // Comments section
    this.commentsElement = insertElement('div', 'gomment-comments', container);

    // Load and render comments
    this.loadCommentsInitial();
  }
}
