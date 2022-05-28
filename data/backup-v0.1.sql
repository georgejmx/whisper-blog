PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE Post (
			id integer primary key autoincrement not null,
			title varchar(40) not null unique,
			author varchar(10),
			contents varchar(1500) not null,
			tag integer not null,
			descriptors varchar(210),
			time datetime default current_timestamp,
			check (tag >= 0 and tag < 8)
		);
INSERT INTO Post VALUES(1, 'Workmates; what?!?', 'pete', 'This blog is all about creating an office community. The key to any career is doing what youre passionate about, which in turn means making strong relationships with colleagues! Events like SupperClub, TheBoardRoom and Henchmen extend the pleasure of work outside the office, so not all about number crunching all day :O. Being lucky enough to work for a great company that invests in its staff is the start; the next piece of the puzzle is valuing colleagues as much as your job. Setbacks at work become areas for improvement, and office conflicts become ripples in a pond, all thanks to the amazing office community I have been lucky enough to join <3',0,'rapid;everlasting;exotic;warmhearted;downright;detailed;square;innocent;fatal;forthright',
'2021-12-02 10:42:26');
INSERT INTO Post VALUES(2, 'The Little Things', 'ataraxy', 'Sitting in the office, listening to the hubbub around me, I learn so much. Ive picked up Python tricks, C# intricacies and other helpful tips that I can use in my everyday work, but these are not the conversations that stick with me. The heated debates on which forms of potatoes are the best (its obviously the daupinoise) and the hilarous recounts of dates gone wrong are the moments that I remember. The obscure facts I pick up from Wikipedia races and the unending discussions about cryptocurrency make up the learning I value most, even if I know that ill probably never use it. After two years of being removed from an office environment, I value these little things the most. The dumb jokes, the awful limericks, putting way too much effort into a set of Christmas decorations to make all the other bays jealous, and rightly so. I may never be a software expert but atleast I know that one day ill master guessing lock screen locations.',3,'alive;overdue;masculine;delightful;acidic;delectable;full;sweet;dependable;thankful',
'2021-12-02 14:07:15');
INSERT INTO Post VALUES(3, 'the wolves', 'musicman', 'rapt in internal cries; oh yee circled by wolves; journey out of the forest; wonder out to the light; grow your strength; see the sky from the heights; the pack will unite',5,'total;keen;gregarious;dreary;imaginative;familiar;each;defiant;unripe;arctic','2021-12-02 15:54:26');
CREATE TABLE Passcode (
			id integer primary key autoincrement not null,
			hash varchar(64) not null
		);
INSERT INTO Passcode VALUES(1,'f480cc5a449100e3a4564127eb1857a607ff06488d2666c254bfa01e6f030bc5');
VALUES(2,'3592db919a6d49f8c79111975aa95299c8ba8e67a807ab0351546dec277649eb');
CREATE TABLE Reaction (
			id integer primary key autoincrement not null,
			postId integer not null,
			descriptor varchar(20) not null,
			gravitas integer not null,
			gravitasHash varchar(64),
			foreign key(postId) references Post(id),
			check (gravitas <= 6)
		);
INSERT INTO Reaction VALUES(1,1,'innocent',2,'');
DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('Post',3);
INSERT INTO sqlite_sequence VALUES('Passcode',2);
INSERT INTO sqlite_sequence VALUES('Reaction',1);
COMMIT;
