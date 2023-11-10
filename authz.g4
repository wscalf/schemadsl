grammar authz;

// Keywords
VERSION: 'version';
VERSIONNO: [0-9]+;
SERVICE: 'service';
ACCESS: PUBLIC | PRIVATE;
PUBLIC: 'public';
PRIVATE: 'private';
ASSET: 'asset';
PERMISSION: 'permission';
LBRACE: '{';
RBRACE: '}';

// Tokens
NAME: [a-z][a-z\-]*;

WS: [ \r\n\t]+ -> skip;

// Rules
file: statement+;

statement: VERSION VERSIONNO #Version
    | SERVICE NAME #Service
    | ACCESS? ASSET NAME LBRACE permission+ RBRACE #Asset
;

permission: ACCESS? PERMISSION NAME;