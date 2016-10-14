import * as React from 'react';

import Chip from 'material-ui/Chip';
import { indigo50 } from 'material-ui/styles/colors';

import { FileMetadata, Highlight } from '../reducers';

const CRLF = /\r\n|\r|\n/g;

interface Props {
    metadata: FileMetadata[];
    contents: Highlight[];
}

interface State {
    highlight: any;
}

const getRawCode = function (reindent) {
    // cached version available ?
    var code = this.rawCode;
    if (code == null) {
        // get the raw content
        code = this.originalCodeblock.get("html");
        // remove empty lines at the beginning+end of the codeblock
        // code = code.replace(/(^\s*\n|\n\s*$)/gi, ""); // @HACK don't repace
        // apply input filter
        code = this.textFilter.filterInput(code);
        // cleanup ampersand ?
        if (this.options.ampersandCleanup === true) {
            code = code.replace(/&amp;/gim, "&");
        }
        // replace html escaped chars
        code = code.replace(/&lt;/gim, "<").replace(/&gt;/gim, ">").replace(/&nbsp;/gim, " ");
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

export class FileContent extends React.Component<Props, State>{

    highlight(offset: number, highlightNums: number[], content: string, ext: string) {
        if (window['EnlighterJS']) {

            // hack
            if (!window['EnlighterJS']._hacked) {
                window['EnlighterJS'].prototype.getRawCode = getRawCode;
                window['EnlighterJS']._hacked = true;
            }

            const pre = document.createElement('pre');
            pre.setAttribute('data-enlighter-lineoffset', String(offset));
            pre.setAttribute('data-enlighter-highlight', highlightNums.join(','));
            pre.appendChild(document.createTextNode(content));

            const div = document.createElement('div');

            const enlighter = new window['EnlighterJS'](pre, {
                indent: 2
            }, div);

            enlighter.enlight(true);

            return <div dangerouslySetInnerHTML={{ __html: div.outerHTML.replace('/\n/g', '') }} />;
        }

        return null;
    }

    render() {
        const { metadata, contents} = this.props;

        contents.sort((a, b) => {
            if ( a.offset < b.offset) {
                return -1;
            }
            if ( a.offset > b.offset) {
                return 1;
            }
            return 0;
        });

        const styles = {
            chip: {
                margin: 4,
            },
            wrapper: {
                display: 'flex',
                flexWrap: 'wrap',
            },
        };
        return (
            <div>
                <div style={styles.wrapper}>
                    {metadata.map(x => {
                        return (
                            <Chip
                                key={`${x.project}_${x.repo}_${x.refs}_${x.path}`}
                                backgroundColor={indigo50}
                                style={styles.chip}
                                >
                                {`${x.project}/${x.repo} (${x.refs}) - ${x.path}`}
                            </Chip>
                        );
                    })}
                </div>
                {contents.map((content, i) => {

                    // stript
                    const codeList = content.content.split(CRLF);

                    // calc highlight lines
                    const highlightNums = codeList.reduce((s, x, index) => {
                        if (x.search(/\u0001/g) !== -1) {
                            s.push(index + content.offset);
                        }
                        return s;
                    }, []);

                    // format
                    const code = codeList.join('\n')
                        .replace(/\u0001/g, '');

                    return (
                        <div key={content.offset}>
                            {this.highlight(content.offset, highlightNums, code, metadata[0].ext)}
                        </div>
                    );
                })}
            </div>
        );
    }
}