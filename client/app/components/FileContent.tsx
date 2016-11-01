import * as React from 'react';

import { Tag } from './Tag';
import { FileMetadata, Preview } from '../reducers';

const CRLF = /\r\n|\r|\n/g;

interface Props {
    metadata: FileMetadata;
    preview: Preview[];
}

interface State {
    highlight: any;
}

export class FileContent extends React.Component<Props, State>{

    highlight(offset: number, highlightNums: number[], content: string, ext: string) {
        if (typeof window !== 'undefined' && window['EnlighterJS']) {

            // hack
            if (!window['EnlighterJS']._hacked) {
                window['EnlighterJS'].prototype.getRawCode = getRawCode;
                window['EnlighterJS'].Renderer.BlockRenderer = CustomBlockRenderer;
                window['EnlighterJS']._hacked = true;
            }

            const pre = document.createElement('pre');
            pre.setAttribute('data-enlighter-lineoffset', String(offset + 1));
            pre.setAttribute('data-enlighter-highlight', highlightNums.map(x => x + 1).join(','));
            pre.appendChild(document.createTextNode(content));

            const div = document.createElement('div');

            const enlighter = new window['EnlighterJS'](pre, {
                indent: 2,
                rawButton: false,
                windowButton: false,
                infoButton: false,
                ampersandCleanup: false
            }, div);

            enlighter.enlight(true);

            return <div dangerouslySetInnerHTML={{ __html: div.outerHTML.replace('/\n/g', '') }} />;
        }

        return null;
    }

    render() {
        const { metadata, preview} = this.props;

        preview.sort((a, b) => {
            if (a.offset < b.offset) {
                return -1;
            }
            if (a.offset > b.offset) {
                return 1;
            }
            return 0;
        });

        const styles = {
            wrapper: {
                float: 'right'
            },
        };
        return (
            <div>
                <div style={styles.wrapper}>
                    {metadata.refs.map(x => {
                        return (
                            <Tag key={x}>{x}</Tag>
                        );
                    })}
                </div>
                <h4>{`${metadata.organization}:${metadata.project}/${metadata.repository} – ${metadata.path}`}</h4>
                {preview.map((pre, i) => {
                    return (
                        <div key={pre.offset}>
                            {this.highlight(pre.offset, pre.hits, pre.preview, metadata.ext)}
                        </div>
                    );
                })}
            </div>
        );
    }
}


// for override Original getRawCode function
const getRawCode = function (reindent) {
    // cached version available ?
    var code = this.rawCode;
    if (code == null) {
        // get the raw content
        code = this.originalCodeblock.get("html");
        // remove empty lines at the beginning+end of the codeblock
        // code = code.replace(/(^\s*\n|\n\s*$)/gi, ""); // @HACK don't remove empty lines
        // apply input filter
        code = this.textFilter.filterInput(code);
        // cleanup ampersand ?
        if (this.options.ampersandCleanup === true) {
            code = code.replace(/&amp;/gim, "&");
        }
        // replace html escaped chars
        // code = code.replace(/&lt;/gim, "<").replace(/&gt;/gim, ">").replace(/&nbsp;/gim, " ");  // @HACK change escape
        code = code.replace(/</gim, "&lt;").replace(/>/gim, "&gt;");

        // @HACK replace gitss tag to <mark> tag
        code = code.replace(/\u0000(.*?)\u0000/g, '<mark>$1</mark>');

        // cache it
        this.rawCode = code;
    }
    // replace tabs with spaces ?
    if (reindent === true) {
        // get indent option value
        var newIndent = this.options.indent.toInt();
        // re-indent code if specified
        if (newIndent > -1) {
            // match all tabs
            code = code.replace(/(\t*)/gim, function (match, p1, offset, string) {
                // replace n tabs with n*newIndent spaces
                return new Array(newIndent * p1.length + 1).join(" ");
            });
        }
    }
    return code;
}

class CustomBlockRenderer {
    options = null;
    textFilter = null;
    constructor(options, textFilter) {
        this.options = options;
        this.textFilter = textFilter;
    }

    render(language, specialLines, localOptions) {
        // elememt shortcut
        var _el = window['EnlighterJS'].Dom.Element;
        // create new outer container element - use ol tag if lineNumbers are enabled. element attribute settings are priorized
        var container = null;
        if (localOptions.lineNumbers != null) {
            container = new _el(localOptions.lineNumbers.toLowerCase() === "true" ? "ol" : "ul");
        } else {
            container = new _el(this.options.showLinenumbers ? "ol" : "ul");
        }
        // add "start" attribute ?
        if ((localOptions.lineNumbers || this.options.showLinenumbers) && localOptions.lineOffset && localOptions.lineOffset.toInt() > 1) {
            container.set("start", localOptions.lineOffset);
        }
        // line number count
        var lineCounter = 1;
        var tokens = language.getTokens();
        var odd = " " + this.options.oddClassname || "";
        var even = " " + this.options.evenClassname || "";
        // current line element
        var currentLine = new _el("li", {
            "class": (specialLines.isSpecialLine(lineCounter) ? "specialline" : "") + odd
        });
        // output filter
        var addFragment = function (className, text) {
            currentLine.grab(new _el("span", {
                "class": className,
                html: this.textFilter.filterOutput(text) // @HACK use 'html' instead of 'text' to write <mark> tag
            }));
        }.bind(this);
        // generate output based on ordered list of tokens
        tokens.forEach(function (token) {
            // split the token into lines
            var lines = token.text.split("\n");
            // linebreaks found ?
            if (lines.length > 1) {
                // just add the first line
                addFragment(token.alias, lines.shift());
                // generate element for each line
                lines.forEach(function (line, lineNumber) {
                    // grab old line into output container
                    container.grab(currentLine);
                    // new line
                    lineCounter++;
                    // create new line, add special line classes; add odd/even classes
                    currentLine = new _el("li", {
                        "class": (specialLines.isSpecialLine(lineCounter) ? "specialline" : "") + (lineCounter % 2 == 0 ? even : odd)
                    });
                    // create new token-element
                    addFragment(token.alias, line);
                });
            } else {
                addFragment(token.alias, token.text);
            }
        });
        // grab last line into container
        container.grab(currentLine);
        // highlight lines ?
        if (this.options.hover && this.options.hover != "NULL") {
            // add hover enable class
            container.addClass(this.options.hover);
        }
        return container;
    }
}
