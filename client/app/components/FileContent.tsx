import * as React from 'react';

import { Grid, Section, Row, Col } from '../components/Grid';
import { RefTag } from './RefTag';
import { FileMetadata, Preview } from '../reducers';

const CRLF = /\r\n|\r|\n/g;

interface Props {
    metadata: FileMetadata;
    keyword: string[];
    preview: Preview[];
}

export class FileContent extends React.PureComponent<Props, void>{

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
                ampersandCleanup: false,
                showLinenumbers: true
            }, div);

            enlighter.enlight(true);

            return <div dangerouslySetInnerHTML={{ __html: div.outerHTML.replace('/\n/g', '') }} />;
        }

        return null;
    }

    render() {
        const { metadata, preview, keyword } = this.props;

        preview.sort((a, b) => {
            if (a.offset < b.offset) {
                return -1;
            }
            if (a.offset > b.offset) {
                return 1;
            }
            return 0;
        });

        keyword.sort((a, b) => {
            if (a.length < b.length) {
                return -1;
            }
            if (a.length > b.length) {
                return 1;
            }
            return 0;
        });
        // console.log(keyword)
        const keywordRegex = keyword.map(x => new RegExp("(" + preg_quote(x) + ")", 'gi'));
        preview.forEach(x => {
            keywordRegex.forEach(re => {
                x.preview = x.preview.replace(re, "\u0000$1\u0000");
            });
        });

        const styles = {
            wrapper: {
                float: 'right'
            },
            title: {
                margin: 0
            },
            row: {
                margin: 0,
                marginBottom: 10
            },
            column: {
                padding: 0
            }
        };
        return (
            <Grid>
                <Row style={styles.row}>
                    <Col size='is12' style={styles.column}>
                        <h4 style={styles.title}>
                            {`${metadata.organization}: ${metadata.project}`}
                            <span style={{ margin: '0 0.25em' }}>/</span>
                            {`${metadata.repository}`}
                        </h4>
                    </Col>
                </Row>
                <Row style={styles.row}>
                    <Col size='is12' style={styles.column}>
                        <h5 style={styles.title}>{`${metadata.path}`}</h5>
                    </Col>
                </Row>
                <Row>
                    <Col size='is10'>
                        {preview.map((pre, i) => {
                            return (
                                <div key={pre.offset}>
                                    {this.highlight(pre.offset, pre.hits, pre.preview, metadata.ext)}
                                </div>
                            );
                        })}
                    </Col>
                    <Col size='is2'>
                        {metadata.branches.map(x => {
                            return (
                                <RefTag key={x} type='branch'>{x}</RefTag>
                            );
                        })}
                        {metadata.tags.map(x => {
                            return (
                                <RefTag key={x} type='tag'>{x}</RefTag>
                            );
                        })}
                    </Col>
                </Row>
            </Grid>
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

// http://stackoverflow.com/questions/280793/case-insensitive-string-replacement-in-javascript
function preg_quote(str) {
    // http://kevin.vanzonneveld.net
    // +   original by: booeyOH
    // +   improved by: Ates Goral (http://magnetiq.com)
    // +   improved by: Kevin van Zonneveld (http://kevin.vanzonneveld.net)
    // +   bugfixed by: Onno Marsman
    // *     example 1: preg_quote("$40");
    // *     returns 1: '\$40'
    // *     example 2: preg_quote("*RRRING* Hello?");
    // *     returns 2: '\*RRRING\* Hello\?'
    // *     example 3: preg_quote("\\.+*?[^]$(){}=!<>|:");
    // *     returns 3: '\\\.\+\*\?\[\^\]\$\(\)\{\}\=\!\<\>\|\:'

    return (str + '').replace(/([\\\.\+\*\?\[\^\]\$\(\)\{\}\=\!\<\>\|\:])/g, "\\$1");
}