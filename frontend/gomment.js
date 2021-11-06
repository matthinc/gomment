/**
 * A response returned by the backend for a comment query.
 * @typedef {Object} CommentQueryResponse
 * @property {Array<Comment>} comments The queried comments.
 * @property {number} total The amount of comments returned by the query.
 */

/**
 * Content of a single comment.
 * @typedef {Object} CommentData
 * @property {number} comment_id Unique identifier for the comment.
 * @property {string} author The author of the comment.
 * @property {number} created_at Creation date of the comment.
 * @property {string} text The content of the comment.
 * @property {number} [num_children] The total number of availabe children, independent of the shown children.
 */

/**
 * A comment response from the backend.
 * @typedef {Object} Comment
 * @property {CommentData} [comment] The actual comment data (content).
 * @property {Array<Comment>} [children] Children comments of this comment.
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

/**
* Transforms a date to a string that will be shown next to the username
* @param {Date} date - the date
* @returns {string}
*/
function defaultDateTransformer(date) {
  return date.toLocaleString();
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
    if(typeof options.thread === 'string') {
      this.threadPath = options.threadPath;
    } else {
      this.threadPath = this.getThreadPathFromLocation();
    }

    // Optional options
    /** @type {number} */
    this.batchSize = options.batchSize || 10;
    /** @type {number} */
    this.maxDepth = options.maxDepth || 2;

    // i18n
    this.i18n = options.i18n || {
      title: 'Comments',
      empty: 'No comments',
      placeholder_name: 'Name',
      placeholder_email: 'E-Mail',
      placeholder_text: 'Your comment',
      submit: 'Submit',
      submit_reply: 'Reply',
      show_more: 'Load more comments',
      show_more_depth: 'Load more replies',
      alert_missing_information: 'Please fill out all required fields!',
      reply: 'Reply',
      new_comment: 'Write comment',
      /** @type {(date: Date) => string} */
      format_date: defaultDateTransformer
    };

    // stateful information
    /** @type {?} */
    this.comments = [];
    /** @type {number} */
    this.lastOffset = 0;
    /** @type {HTMLElement | null} */
    this.commentsElement = null;
    /** @type {HTMLElement | null} */
    this.submitButton = null;
    /** @type {HTMLElement | null} */
    this.replyIndicator = null;
    /** @type {HTMLElement | null} */
    this.newButton = null;
  }

  /**
   * Query for comments with the specified parameters.
   * @param {number} offset - The offset.
   * @param {number} max - Maximum amount of comments to query.
   * @param {number} depth - Maximum level of depth to query.
   * @returns {Promise<Response>} - HTML Response.
   */
  queryComments(offset, max, depth) {
    return window.fetch(`${this.apiURL}comments?threadPath=${encodeURIComponent(this.threadPath)}&offset=${offset}&max=${max}&depth=${depth}`);
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

    this.queryComments(offset, max, depth)
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
   * Reload comments with current parameters
   * @returns {void}
   */
  reloadComments() {
    this.loadComments(this.lastOffset, this.batchSize, this.maxDepth);
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
    let commentClass = 'gomment-comment';
    if (comment.comment.comment_id === this.newComment) {
      commentClass += ' gomment-comment-new';
    }

    const commentElement = insertElement('div', commentClass, parent);
    insertElement('div', 'gomment-comment-author', commentElement, { innerHTML: comment.comment.author });
    insertElement('div', 'gomment-comment-date', commentElement, { innerHTML: this.i18n.format_date(new Date(comment.comment.created_at * 1000)) });
    insertElement('div', 'gomment-comment-text', commentElement, { innerHTML: comment.comment.text });

    const replyButton = insertElement('a', 'gomment-comment-reply', commentElement, { innerHTML: this.i18n.reply });

    const childrenElement = insertElement('div', 'gomment-comment-children', commentElement);

    replyButton.onclick = () => {
      this.moveCommentTarget(childrenElement, comment.comment.comment_id);
    };

    if (comment.children) {
      for (const childComment of comment.children) {
        this.renderComment(childrenElement, childComment, treeIndex, depth + 1);
      }
    }
    if (comment.comment.num_children > comment.children.length) {
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
   * Publish new comment to this thread
   * @param {number} parent - parent comment (0 = top-level)
   * @param {string} name
   * @param {string} email
   * @param {string} content - text content
   * @returns {boolean} true if data was valid
   */
  publishComment(parent, name, email, content) {
    if (name && email && content) {
      const data = {
        author: name,
        email,
        text: content,
        parent_id: parent,
        thread_path: this.threadPath,
      };

      // Post comment
      window.fetch(`${this.apiURL}comment`, {
        method: 'POST',
        body: JSON.stringify(data)
      })
        .then(response => response.json())
        .then(data => {
          this.newComment = data.id;
          this.reloadComments();
        });

    } else {
      alert(this.i18n.alert_missing_information);
      return false;
    }

    return true;
  }

  /**
   * Move comment input to new parent and change recipient
   * @param {HTMLElement} parent - new parent
   * @param {number} recipient - new recipient (0 for top-level comments)
   */
  moveCommentTarget(parent, recipient) {
    // Set recipient
    this.replyRecipient = recipient;

    // Create new
    if (!this.inputSection) {
      const inputSection = insertElement('div', 'gomment-input-section', parent);
      const mailElement = /** @type {HTMLInputElement} */ (insertElement('input', 'gomment-email', inputSection, { placeholder: this.i18n.placeholder_email }));
      const nameElement = /** @type {HTMLInputElement} */ (insertElement('input', 'gomment-display-name', inputSection, { placeholder: this.i18n.placeholder_name }));
      const contentElement = /** @type {HTMLInputElement} */ (insertElement('textarea', 'gomment-text-input', inputSection, { placeholder: this.i18n.placeholder_text }));

      const publish = () => {
        const recipient = this.replyRecipient || 0;
        if (this.publishComment(recipient, nameElement.value, mailElement.value, contentElement.value)) {
          // Clear inputs
          mailElement.value = '';
          nameElement.value = '';
          contentElement.value = '';
        }
      };

      insertElement('button', 'gomment-submit-button', inputSection, { innerHTML: this.i18n.submit, onclick: publish});

      // Cache
      this.inputSection = inputSection;
    }

    // Move
    if (parent.childNodes.length === 0) {
      parent.appendChild(this.inputSection);
    } else {
      parent.insertBefore(this.inputSection, parent.childNodes[0]);
    }

    // Hide new comment button
    if (this.newCommentButton) {
      this.newCommentButton.hidden = recipient === 0;
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

    // New comment button
    const inputContainer = insertElement('div', 'gomment-new-comment-input-container', container);

    this.newCommentButton = insertElement('button', 'gomment-new-comment', inputContainer, { innerHTML: this.i18n.new_comment});

    this.newCommentButton.onclick = () => {
      this.moveCommentTarget(inputContainer, 0);
    };

    // Initial input field position
    this.moveCommentTarget(inputContainer, 0);

    // Comments section
    this.commentsElement = insertElement('div', 'gomment-comments', container);

    // Load and render comments
    this.loadCommentsInitial();
  }

  /**
   * Calculate a cleaned-up thread path based on the current browser
   * location by trimming all trailing slashes.
   * @returns {string} - The cleaned up thread path
   */
  getThreadPathFromLocation() {
    let pathname = window.location.pathname;

    // remove trailing slash if it exists
    if (pathname[pathname.length - 1] === '/') {
      pathname = pathname.substring(0, pathname.length - 1);
    }

    // if the path is empty, use a slash as path
    if (pathname.length === 0) {
      pathname = '/';
    }

    return pathname;
  }
}
