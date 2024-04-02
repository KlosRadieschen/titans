package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
)

var (
	sessions        [4]*discordgo.Session
	personalities   []*discordgo.Webhook
	awaitUsers      []string
	awaitUsersDec   []string
	missionUsers    []string
	missionChannels []string
	donator         string
	donatorRole     string
	sacrificed      bool
)

var (
	GuildID  = "1195135473006420048"
	sleeping = []bool{false, false, false, false}
	modes    = make(map[string]bool)
	message  = make(map[string][]string)
	client   = openai.NewClient(openAIToken)
	req      = openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `
Here is the full database of the AHA (Anti-Horny Alliance) discord server:

CREATE TABLE Fleet (pkfk_battalion_owns TINYINT unsigned NOT NULL,carriers SMALLINT unsigned NOT NULL,battleships SMALLINT unsigned NOT NULL,heavyCruisers SMALLINT unsigned NOT NULL,lightCruisers SMALLINT unsigned NOT NULL,destroyers SMALLINT unsigned NOT NULL,frigates SMALLINT unsigned NOT NULL,dropships SMALLINT unsigned NOT NULL,transportships SMALLINT unsigned NOT NULL,fk_flagship_leads VARCHAR(32) NOT NULL,PRIMARY KEY ('pkfk_battalion_owns'));
INSERT INTO Fleet VALUES(4,3,2,5,10,15,30,240,3,'Meruda');
INSERT INTO Fleet VALUES(2,20,2,6,6,3,15,200,0,'Midas');
INSERT INTO Fleet VALUES(0,0,0,0,0,0,2,0,0,'Infiltrator');
INSERT INTO Fleet VALUES(3,10,3,13,16,29,57,0,0,'Rift Walker');
INSERT INTO Fleet VALUES(1,36,10,28,18,52,70,540,0,'Resolute');
INSERT INTO Fleet VALUES(5,20,5,18,10,40,60,300,0,'Vanguard');
CREATE TABLE Flagship (pk_name VARCHAR(32) NOT NULL,class VARCHAR(32) NOT NULL,description TEXT,titanCapacity INT NOT NULL,PRIMARY KEY (pk_name));
INSERT INTO Flagship VALUES('Meruda','Dreadnought','',100);
INSERT INTO Flagship VALUES('Midas','Dreadnought','',100);
INSERT INTO Flagship VALUES('Infiltrator','Light Cruiser','',6);
INSERT INTO Flagship VALUES('Resolute','Dreadnought','Gigantic ship carrying a fold weapon',32);
INSERT INTO Flagship VALUES('Rift Walker','Dreadnought','',100);
INSERT INTO Flagship VALUES('Vanguard','Dreadnought','It is a ship that can fly',100);
CREATE TABLE Base (pk_name VARCHAR(20) NOT NULL,size VARCHAR(20),fk_planet_isOn VARCHAR(20) NOT NULL,PRIMARY KEY (pk_name));
INSERT INTO Base VALUES('Gaia','Large outpost','Laythe');
INSERT INTO Base VALUES('Novus','Capital','Typhon');
INSERT INTO Base VALUES('Halcyon','AHA Main Base','Harmony');
CREATE TABLE Planet (pk_name VARCHAR(20) NOT NULL,environment TINYTEXT,fk_system_isInside VARCHAR(20),fk_battalion_controls TINYINT NOT NULL,PRIMARY KEY (pk_name));
INSERT INTO Planet VALUES('Laythe','A relatively small planet with a lot of oceans. However, the atmosphere is breathable for humans and both flora and fauna are friendly','Jool',4);
INSERT INTO Planet VALUES('Harmony',replace('Harmony is a planet within the Freeport System of the Frontier. It is an agricultural world and the HQ of the A.H.A. The planet has at least one moon, and a population of 40 million and is run by the 1st Battalion. \nIt is a peaceful and calm planet, with sleek stone and wooden architecture, along with plenty of hexagonal rock formations. As a primarily agricultural world, its landscape is littered with agricultural towers and fields of crops. \nAs Q the A.H.A, Harmony is a heavily defended world, being home to the Harmony Defensive Fleet.','\n',char(10)),'Freeport',1);
INSERT INTO Planet VALUES('Typhon','','Sector Echo-152',2);
INSERT INTO Planet VALUES('Orthros','','Sector Echo-152',3);
CREATE TABLE System (pk_name VARCHAR(32) NOT NULL,PRIMARY KEY (pk_name));
INSERT INTO System VALUES('Jool');
INSERT INTO System VALUES('Freeport');
CREATE TABLE Battalion (pk_number TINYINT UNSIGNED NOT NULL, name VARCHAR(32), fk_pilot_leads VARCHAR(18) NOT NULL, PRIMARY KEY (pk_number));
INSERT INTO Battalion VALUES(4,'Phoenixes of Laythe','384422339393355786');
INSERT INTO Battalion VALUES(0,'SWAG','455833801638281216');
INSERT INTO Battalion VALUES(1,'Praetorians of Harmony','1079774043684745267');
INSERT INTO Battalion VALUES(2,'Typhon Grenadiers','992141217351618591');
INSERT INTO Battalion VALUES(3,'Heh','1022882533500797118');
INSERT INTO Battalion VALUES(5,'Defense of Harmony','1079774043684745267');
CREATE TABLE PersonalShip(pk_name VARCHAR(32), class VARCHAR(20) NOT NULL, description TEXT, titanCapacity TINYINT, PRIMARY KEY (pk_name));
INSERT INTO PersonalShip VALUES('Reaper','Large Fighter','A figher relatively fast for it''s size. Advanced stealth makes it impossible to detect on most radars.',1);
INSERT INTO PersonalShip VALUES('Castellum','MacAllan (carrier)','Heavily fortified and slightly modified. Carries many bombers, fighters and dropships with a full complement of crew and grunts','16 vanguard class');
INSERT INTO PersonalShip VALUES('Mystery','Recon','a small fast ship able to house 10 pilots and 2 titans',2);
INSERT INTO PersonalShip VALUES('Flop-pod','Drop-Pod','A custom drop pod with superb manuevering and coffee machine',0);
INSERT INTO PersonalShip VALUES('AHF Rift Walker','Super Carrier','Being able to tear rifts to enter P-Space, Rift Walker can get in and out of combat quickly if needed.',15);
INSERT INTO PersonalShip VALUES('Star of Steel','Star Destroyer','It can fabricate items, airdrop them, send big bombs and wipe out enemies',4);
INSERT INTO PersonalShip VALUES('Star Jumper II','Carrier-Assoult-Extraction','3 turrets [1 top turret, 2 side turrets] - Camouflage ability can make you disappear from radar - 4 engine dark propulsion [top speed: mack 20] [Hyperspace cooldown 3 sec] - The ship can be used for close-high-risk missions and can allow rapid engagement of the target and rapid extraction even for Titans thanks to its tailgate','6 [4 pod titanfall]');
INSERT INTO PersonalShip VALUES('castigator','destroyer','large yet fast warship ready to kick ass',20);
INSERT INTO PersonalShip VALUES('JumperII','Carrier-Assoult-Extraction','3 turrets [1 top turret, 2 side turrets] - Camouflage ability can make you disappear from radar - 4 engine dark propulsion [top speed: mack 20] [Hyperspace cooldown 3 sec] - The ship can be used for close-high-risk missions and can allow rapid engagement of the target and rapid extraction even for Titans thanks to its tailgate','6 [4 pod titanfall]');
INSERT INTO PersonalShip VALUES('Jumper','Carrier','3 turrets [1 top turret, 2 side turrets] - 4 pods for Titanfall - Camouflage ability can make you disappear from radar - 4 engine dark propulsion [top speed: mack 20] [Hyperspace cooldown 3 sec] - The ship can be used for close-high-risk missions and can allow rapid engagement of the target and rapid extraction even for Titans thanks to its tailgate',6);
INSERT INTO PersonalShip VALUES('Survivor','corvette','a large corvette boarding on frigate class. It is designed for ramming and board opponents ships. It has a powerful boardsides to beat any ship of the class. It is a quick ship with weaker shields and weak against aircraft. It has trouble staying in atmosphere for long periods making it poor for accurate ground support.',5);
INSERT INTO PersonalShip VALUES('survivor','corvette','The corvette has a pair of twin 300inch gun batteries on each side along with many AA guns being replaced for boarding ramps as the ship is designed to board enemy ships to quickly disable them. Along with that it has a heavy blade on the front for ramming enemy vessels the disadvantage of this meaning the ship can’t stay in atmosphere for long periods of time. Though marked as corvette its size can make it rival some frigates. It is well able take most ships in it weight class',15);
INSERT INTO PersonalShip VALUES('Executioner','cruiser','big railgun and ailien tech',75);
INSERT INTO PersonalShip VALUES('Resolute','Dreadnought',replace('Flagship of the 1st Fleet and the A.H.A as a whole. \n\nIt is the largest ship owned by the A.H.A sporting extremely heavy armour and weaponry. \n\nIt utilizes the Ark as a power source with the device being built into the ships centre it can also utilize it to fire the devastating fold weapon capable of destroying any enemy vessel. \n\nWhilst lacking the carrying capacity of some of the other dreadnoughts it makes up for it in ship to ship combat abilities.','\n',char(10)),32);
CREATE TABLE Titan(pk_callsign VARCHAR(7), name VARCHAR(32) , class VARCHAR(20) NOT NULL, weapons TINYTEXT, abilities TINYTEXT, PRIMARY KEY (pk_callsign));
INSERT INTO Titan VALUES('VN-8385','Vivian','Vanguard atlas/ogre','None','weapon adaptation, vortex shield and particle wall');
INSERT INTO Titan VALUES('MW-0451','Espresso','Ogre','Combustion Pipe, Full-Auto Thermite Launcher','MicroWave Power Overload, Coffee Machine');
INSERT INTO Titan VALUES('TK-4804','Big Red','Legion','Predator Cannon','Power Shot, Gun Shield, Long Range Mode, Smart Core and Advanced AI');
INSERT INTO Titan VALUES('CD-1050','Charlie','Northstar','Plasma Railgun, cluster missiles, rockets','Flight, Nuclear Core');
INSERT INTO Titan VALUES('Hangman','Butcher','Vanguard','sword XO17','vortex shield tracking rockets');
INSERT INTO Titan VALUES('HB-8455','Jorge','Legion','Predator cannon','Smart core, gun shield, speed mod, missiles (thanks Gold)');
INSERT INTO Titan VALUES('QL-9821','henry','Recon , Northstar','Arc cannon and Plasma railgun','able to fly like viper');
INSERT INTO Titan VALUES('YB-2784','Cornelius','stryder','heavy plasma railgun','plasma EMP shot, enhanced flight core, vortex shield, missile salvo, and thrusters');
INSERT INTO Titan VALUES('TR-5834','Tempest','Stryder','PR-01 Plasma Railgun','Advanced Hover, Cluster Salvo, Particle Wall');
INSERT INTO Titan VALUES('THE-J','The J','J class','J, air support marker, missile launcher, white phosphorus launcher, ULTRATURRET','jump, nuke eject, breakdance, trow it back');
INSERT INTO Titan VALUES('idk','crabegion','ogre','predator cannon','flight');
INSERT INTO Titan VALUES('DS-2629','Jeff','Scorch','Wildfire Launcher, Warcrime Traps','Firewall, Dementia');
INSERT INTO Titan VALUES('AO-0740','Alpha','Tone','44 mm tracker cannon, scythe','Reinforced Plasma Wall, Sonar Pulse');
INSERT INTO Titan VALUES('CR-2007','Commissar’s retribution','tone','modified 88mm AT cannon that can be used dual with a modified sword (chainsword)','it has a armour more like a legion and has improved sensors. It can also phrase dash along with all normal tone abilities. Its main pistol can be supercharged to increase penetration power. Chainsword is effective against infantry along with a second mode used against titans');
INSERT INTO Titan VALUES('AT-8082','hornet','atlas','40 mm tracker cannon','tracking missiles, reflective particle wall, sonar pulse');
CREATE TABLE Pilot (pk_userID VARCHAR(18) NOT NULL, platform VARCHAR(40), ingameName VARCHAR(32), specialisation VARCHAR(50), isSimulacrum BOOLEAN NOT NULL, story TEXT, fk_battalion_isPartOf TINYINT unsigned,fk_personalShip_possesses VARCHAR(32), fk_titan_pilots VARCHAR(7) , PRIMARY KEY (pk_userID));
INSERT INTO Pilot VALUES('455833801638281216',NULL,NULL,'SWAG leader',1,NULL,0,NULL,'CD-1050');
INSERT INTO Pilot VALUES('1022882533500797118',NULL,NULL,'Weapon Seller',1,NULL,3,'AHF Rift Walker','TR-5834');
INSERT INTO Pilot VALUES('985218545677922375',NULL,NULL,'Recon Specialist',0,'I was drafted into the Militia at 18 years old. Born and raised on Harmony, life used to be tough for me and my family. When I got drafted, my family was left without anyone to help everyone survive. Chances are, they''re dead now. After training in the Militia, I got assigned to a squad. ''Solstice Dawn'', they were called. It was good. I stayed with them for 2 years, until... one night, my robotic leg got hacked. By the time I woke up and stopped it, the last person in the squad, my commanding officer, was gone. I... don''t like talking about it. I started heading back to base, when I found a Titan. Callsign Alpha Omega 740, from the Marauder Corp. He had just lost his Pilot, and offered to take me back to base. We bonded, and on arriving to base, we got officially linked. That was on the third day of me knowing him. Not even a month later, I got discharged from service. Dishonourable discharge. Alpha came with, since we were linked. I wouldn''t leave him. Soon after, PHC tried drafting me. I saw the AHA and approved of it. So I joined when they reached out. Now I''m here, and I''ve found... at least some semblance of a place to belong.',4,'','AO-0740');
INSERT INTO Pilot VALUES('384422339393355786','PC','KlosRadieschen','Chief Hackerman',1,'I worked at Hammond Robotics for a long time and this is where I learned most of what I know today. As you can probably imagine, I helped researching modern simulacrum technology but I also contributed to the code of titan AIs.  When the IMC suffered losses in the war with the militia, they started searching more men, even going to their allies and taking people that were not originally soldiers and I was also chosen.  I technically wouldn''t even get a titan but my former colleagues and I build a Scorch for me and had to do all of the coding ourselves.  I fought a few battles but I really felt like the IMC didn''t appreciate my work at all and I decided to leave it all behind, taking only my titan with me. I was lost for some time until I discovered the AHA. At first, I only heard of the war but after giving it some thought, I decided to join and fight with my Scorch.  However, after some minor battles, my Scorch started developing dementia because of the code we wrote at Hammond and I wasn''t sure about fighting anymore. So I did a step back into my past and decided to contribute to the AHA in a different way.',4,'Reaper','DS-2629');
INSERT INTO Pilot VALUES('1054598086191759372',NULL,NULL,'Excessive Force',0,NULL,2,'','CF-4952');
INSERT INTO Pilot VALUES('952145898824138792',NULL,NULL,'medic',1,'born on a  colonised outer world into a slave market. At the age of 15 escaped however leaving my mother and sister behind. Found refuge as a rifleman before training to become a pilot. Joined the 6-4 and stayed there until the age of 24 when on an assassination mission was killed. Found by IMC scientists, I was made into a unwavering killing machine, murdering an entire planet containing my entire 6-4 battalion stationed on there before getting back control of my own body. Needing anything to keep my alive and give me a purpose, I joined up with the apex predators for mercenary work at the age of 46. I was eventually brought into that of the horny wars as a general under the command of Ender_Fender. Still a devout AHA member and still clinging to mercenary work as it''s my only sence I have left of my humanity....',4,'castigator','idk');
INSERT INTO Pilot VALUES('780847814753779770',NULL,NULL,'SwordMedic',1,'Fish was a guy. He became the second commissar of our faction and fought side-by-side with our men in many missions from all battalions. But then, the incident happened. In the Defense of Laythe battle, a weird bug-like microbot was able to completely take over Fish''s body, which had just become a simulacrum. After dropping in his titan through the roof of the 4th battalion flagship, his possessed body almost killed his most handsome, intelligent and definitely most modest friend, Commander Klos. After having to witness that, his mind snapped. He became what can be best described as Spiderman possessed by a crackhead demon. However, this is also the time he started doing some incredible and complex works of engineering.',5,'Executioner','Hangman');
INSERT INTO Pilot VALUES('942159289836011591',NULL,NULL,'Combat medic and legion pilot',0,NULL,0,'whatever is the ship of swag','VQ-9823');
INSERT INTO Pilot VALUES('920342100468436993',NULL,NULL,'Ships and equipment',0,replace('I was enlisted in an IMC military academy from a very young age at the behest of my father (Captain [REDACTED]). I was top of my class and eventually graduated and was given 2nd in command of an IMC battleship. During an operation near the [REDACTED] system, the captain of my ship was killed. I had to take command but lost almost all of my men and all secondary IMC ships. During this time my father was killed defending a key IMC manufactorum world, this meant he couldn’t vouch for me when I was dishonourably discharged for gross incompetence. They needed someone to blame and (since the captain was dead) it had to be me.\n\n After this I fell into mercenary work worth a group known as [REDACTED] for 3 years before my ship and crew were killed in an incident that…I don’t want to go into. I found the AHA, a cause I finally believed in, and became a Lieutenant (my mercenary group weren’t part of the centralised war effort but still fought the horny). \n\nI climbed to the rank of Major and then the 2nd battle of Laythe happened. I was critically injured by the impact of a capital ship and the medics aboard the Midas were revealed to be endless spies. They implanted me with some sort of alien parasite. I almost destroyed myself to terminate it and needed both my arms replaced, as well as part of my skull. Since then I was promoted to Colonel and haven’t been promoted since as there are no positions in high command free.','\n',char(10)),2,'Castellum','HB-8455');
INSERT INTO Pilot VALUES('1079774043684745267',NULL,NULL,'Leader of the A.H.A',0,NULL,1,'Resolute','TK-4804');
INSERT INTO Pilot VALUES('992141217351618591',NULL,NULL,'Nukes',1,NULL,2,'','TU-8271');
INSERT INTO Pilot VALUES('1136401069908426894',NULL,NULL,'Multi-tasker, jack and master of all trades',1,NULL,2,'Mystery','QL-9821');
INSERT INTO Pilot VALUES('629044798799478815',NULL,NULL,'engineering',1,NULL,2,'saber-1 and saber-2','YB-2784');
INSERT INTO Pilot VALUES('634554969814597667',NULL,NULL,'cooking and modifying weapons and building titans',0,NULL,5,'','VN-8385');
INSERT INTO Pilot VALUES('707016988140240926',NULL,NULL,'Engineering, field tests, vehicle/titan piloting',0,NULL,2,'Flop-pod','MW-0451');
INSERT INTO Pilot VALUES('947109747390300210',NULL,NULL,'Hydrogen Combustion',1,'I was born in Portugal, overthrew the goverment, invaded every neighbouring country, lost, got the death sentence, left earth, found myself a good place to stay, enslaved the natives, started researching nanobot technology, the natives revolted, left, found a place where they were fighting with robots and jetpacks, started researching it, made contact with a group called W.A.S.P., Commisar Raptor contacted me and invited me to the A.H.A. I accepted, fought in the 2nd horny war, joined the swag, found a great enthusiasm in bombing hospitals, got cancer causing nanobots, had to turn myself into a simulacrum, i forgor 💀, I sent myself into space, enslaved natives on a planet, made ship and robot factories, started using robots to talk (and bomb) these people, had a fleet of 100,000 ships, gave up on that, gave up on the space thing, realised I was doing all the messed up stuff because I was searching for something that doesn''t exist and now I''m trying to redeem myself.',0,'Star of Steel','THE-J');
INSERT INTO Pilot VALUES('1138968064851972157',NULL,NULL,'Weapon Specialist, Recon',0,NULL,2,'Star Jumper II','CA-2212');
INSERT INTO Pilot VALUES('1050395672366555156',NULL,NULL,'not dieing',0,NULL,4,'survivor','CR-2007');
INSERT INTO Pilot VALUES('554150090164535322',NULL,NULL,'Combat Medic',0,replace('Back when I was a kid, the PHA visited my little part of Cibus.\n\nThey told us they weren’t staying long, just enough to recoup themselves and all that. They ended up staying for a while, preaching about loving openly and all that good propaganda crap\n\nThe one big promise they kept for a while was that the war would never come to us. We were adamant that we wanted nothing to do with the first war. They listened for a while, until they stuck a superweapon basically in my backyard and called it an Art Project\n\nLike usual we all listened to their lies, until they fired the weapon and we saw through all the fake lies and shark tooth smiles they gave us for years\n\nI was angry, confused, and had gotten separated from my home and family in the middle of the conflict. Either I could join the AHA, or be more human debris floating through space.\n\nI made my decision, and I ain’t lookin back','\n',char(10)),2,'','Providence');
INSERT INTO Pilot VALUES('462970583122968587',NULL,NULL,'forward assault group',0,NULL,1279,'','Archangel');
INSERT INTO Pilot VALUES('989615855472082994',NULL,NULL,'Grunt',0,NULL,4,'','none');
INSERT INTO Pilot VALUES('754699286440050708',NULL,NULL,'none',0,NULL,4,'','AT-8082');
CREATE TABLE Report (pk_name VARCHAR(64) NOT NULL,timeIndex SMALLINT NOT NULL,type TINYINT NOT NULL,authorType TINYINT NOT NULL,fk_pilot_wrote VARCHAR(18) NOT NULL,description TEXT NOT NULL,PRIMARY KEY (timeIndex));
INSERT INTO Report VALUES('Battle of the Last Tesseract',0,3,1,'992141217351618591',replace('Krazy and Fish went to scout out for the final Tesseract. They requested forces on standby in case they needed them.\n\nThey found the Tesseract and called backup. \n\nMe and fish retrieved the memory core. We now have everything about the Endless.\n\nThey called in a fleet and the Midas started taking damage.\n\nI evacuated the forces that had helped, including the swag and 4th battalion who aided in destroying the enemy fleet, and then detonated a nuclear warhead, killing the enemies.\n\nWe had minimal injuries and sustained no deaths.\nNo friendly forces were caught in the blast.','\n',char(10)));
INSERT INTO Report VALUES('Reconquest of Typhon & Rooster''s warcrimes',-80,3,1,'989615855472082994',replace('Com. Rooster. The battle for Typhon. \nOnce the vengeance was completed we moved out to retake typhon. TU had idiotically went in alone, and everyone else was dispatched to reinforce. I still can''t believe what I''ve done, they were surrendering and I didn''t stop. I said no survivors and continued slaughtering them as they surrendered, I have a feeling many of the others felt the way I did. Because no one stopped me. Not when I commanded the execution of those that surrendered, no I did anyone try to stop me when I called in the life-eater. I underestimated the blast radius. The capital was caught. 700,000 people, dead, because of me, afterwards me and TU went down to look for survivors. Six, only six survivors out of 700,000. Now we wait for the general to find out, he will decide what to do with me. And frankly I hope he chooses death. If he doesn''t I''ll do it for him.\nCom. Rooster signing off, for the final time.','\n',char(10)));
INSERT INTO Report VALUES('Death of JyFK',-50,5,1,'1079774043684745267',replace('Name : JyFK\n\nRank : Commander\n\nStatus : K.I.A\n\nCommander JyFK served honourable since the end of the First War, however, overtime his mental state began to deteriorate. It is unclear if this was a side effect of his transfer to a simulacrum body or just the stresses of his command. \n\nIt was gradual at first but the death of the former Commander Rooster, who had once been JyFKs second in command and dear friend pushed him over the edge. Adding to this further was that Rooster took his life with the gun he had taken from JyFKs person during their final dialogue. \n\nThe guilt and grief proved too much for Commander JyFK and he went AWOL. He attempted a suicide mission on the enemy commandeering a ship and severely wounding several friendly personnel almost killing some before departing. \n\nHe continued to fire upon enemy and friendly units indiscriminately, the 1st Battalion was deployed with a battlegroup jumping in on his location at which point from the lead ship of the battlegroup I activated a large scale EMP. \n\nThis prevented the now inorganic Commander from moving and locked his consciousness to one body, from there his ship was boorded, a restraining bolt fitted and the Commander was taken into custody. \n\nHe was questioned by several lower ranking members as well as myself whilst in a temporary holding cell on Harmony where he was awaiting trial. \n\nIt was then that he activated an explosive device in his chest in an attempt to kill everyone present, I promptly shot the Commander in the head with my wingman in an attempt to deactivate the device but to no avail. \n\nThe device exploded injuring several personnel who were rushed to the medbay and luckily were stabilized by medical personnel resulting in no casualties except the Commander himself. \n\nSadly his final series actions tainted his other wise exemplary service record leaving a dark stain on his legacy.\n\nThis will not be allowed to happen again. \n\nLong Live the A.H.A','\n',char(10)));
INSERT INTO Report VALUES('Death of Rooster',-70,5,1,'1079774043684745267',replace('Name : Rooster\n\nRank : Commander\n\nStatus : K.I.A\n\nFollowing the failed operation on Typhon resulting in the deaths of 700,000 civilians and destruction of the capital Commander Rooster was apprehended and taken into Custody to await trial. \n\nThe trial was broadcast and the Commander was stripped of his command and arrested. \n\nCommander JyFK wished to say his goodbyes so I escorted him to the holding facility. We both entered the cell and had a brief dialogue with the former Commander. Unbeknownst to us the former Commander had stolen Commander JyFKs sidearm during the dialogue. \n\nUpon our departure we stopped briefly in the hallway outside the holdings cells room to discuss the events that had transpired. This was when we heard the gunshot. \n\nOvercome with guilt for his actions and anxiety over the impact his actions would have on our forces the former Commander used the stolen sidearm to take his own life. \n\nDespite the best efforts of nearby medics he could not be saved and passed away. \n\nHe will be missed. \n\nLong Live the A.H.A\n','\n',char(10)));
INSERT INTO Report VALUES('Death of John',10,5,1,'1079774043684745267',replace('Name : John\n\nRank : Rifleman\n\nStatus : K.I.A.\n\nRifleman John, brother to the late Commander JyFK was a promising soldier. \n\nHowever, upon denial of his request to revive his older brother he began to grow distant. \n\nThis resulted in John going AWOL, he shot and wounded a dropship Pilot before stealing said dropship and attempting to reach his brothers burial site in an effort of bringing him back. \n\nUnder direct orders from myself,  A.H.A special forces (S.W.A.G) were deployed to handle the situation. \n\nTheir light cruiser the A.H.F infiltrator performed a jump coming out in close proximity to the stolen vessel at which point he was ordered to stand down or be fired upon. \n\nHe chose the latter, fighter craft were deployed dealing critical damage to the dropship causing it to go down. It crash landed near an urban area on a nearby planets surface. \n\nSeveral S.W.A.G pilots were deployed to the ground to capture the AWOL rifleman. \n\nWeather was bad with torrential downpour of rain leading to low visibility which ultimately aided our operatives. \n\nAfter stalking the rifleman from a distance and a brief chase he was cornered on the rooftop of one of the crowded cities buildings, where ultimately he was shot and killed.\n\nThis was after resisting arrest and attempting to draw a gun on the operatives. \n\nHis body laid there briefly in the downpour of rain, slowly staining the puddle he crumpled in as the weather worsened and storms began to form. \n\nHis body was recovered and his being moved for burial. ','\n',char(10)));
INSERT INTO Report VALUES('Destruction of endless planet',20,3,1,'384422339393355786','The 2nd battalion detected a hostile planet not very far from Laythe. High command made the decision to attack that planet with the combined force of the entire AHA. The second battalion arrived first with the others only seconds after and we started an assault on the ground. We advanced quickly towards the the main base. I (Klos) was able to detect enemy backup early, and defeated them before they could group up with the main forces thanks to the help of Cipher. Jesse and Phantom were hurt while fighting in the base. Eventually, TU set up the explosion of the entire planet and we evacuated successfully, with no further casualties.');
INSERT INTO Report VALUES('Death of Ender',30,5,1,'1079774043684745267',replace('Betrayal seemed to cling to Ender like a foul stench, following him wherever he goes. \n\nHe founded our organisation, he led us against seemingly insurmountable odds and I followed his every order with admiration. But that was a long time ago. \n\nOnce again he has brought death to our door, endangering our forces.. My friends and brotherd. \n\nThis is but one final stain on his once glorious story, he fell from grace. Whilst chaos erupted in every direction, the muffled sound of gunfire and explosions filling the air. \n\nYet that room was so quiet. \n\nHe was fitted with a restraining bolt and in that moment he was just as mortal as any of us. It is in that moment I executed him for his crimes, for his betrayal. \n\nYou started a great organisation Ender, when you get to hell.. \n\nSave me a seat. \n\nLong Live the A.H.A','\n',char(10)));
INSERT INTO Report VALUES('Enders Removal from Office',-240,4,1,'1079774043684745267',replace('I feel this war is coming to an end. Yet with victory in our reach a vile truth came to light, it was revealed that our general.. Has betrayed us. \n\nGeneral Ender, codenamed Viper has fallen to the enemies tainted ideologies and forsaken everything we have fought and died for. He has sold us out to the enemy and endangered all our lives. \n\nI held a vote amongst our most senior officers including my fellow commanders, captains and lieutenants. \n\nThe decision was made. \n\nI took a squad of my most elite men and personally stormed the facility Ender was stationed in. We were unopposed for the most part once Enders crimes came to light, there were some devote fanatics however, although they were dealt with quickly. \n\nWe breached the room and took Ender into custody. \n\nAs of this moment. I am in control. I will not fall victim to the pitfalls of my predecessor.\n\nLong Live the A.H.A','\n',char(10)));
INSERT INTO Report VALUES('A.H.A Civil War',-250,3,1,'1079774043684745267',replace('I have been chief warden of Enders containment facility for some time now, the war still rages and though I''ve been deployed more times than I can count this is my main post between deployments. \n\nA few weeks back a subfaction of our forces formed a union because they were disgruntled over lack of pay.. How selfish. \n\nThis war is bigger than all of us, the honour of upholding my duty is pay enough for me. If you ask me they were all traitors, all they are doing is aiding the enemy by weakening our forces and dividing our attention. \n\nThis culminated in them mobilising against us, after several strikes they moved to assault Enders main containment facility, where I was stationed. \n\nI met them in the planes in front of the facility, I issued them with several orders to stand down. Though this was in vain, this deadlock went on for several hours then a single shot was fired. It was unclear which side had fired but it was the spark they were waiting for. \n\nA skirmish broke out with both sides opening fire, we were pushed back to the gate but we gave up no more ground than that. I cut down many of my supposed allies, we would not yield. \n\nDuring the chaos a P.H.C infiltrator managed to gain access to the facility freeing the prisoners. I killed many prisoners and traitors alike though despite my best efforts many prisoners escaped. \n\nAfter the massacre at the facility the union was weakened significantly and later came to an agreement with Ender. \n\nIt was truly a dark day. \n\nLong Live the A.H.A','\n',char(10)));
INSERT INTO Report VALUES('The Incident',-20,6,1,'384422339393355786','The incident happened just after the "Defense of Laythe" battle, in which we found out that the endless are capable of taking over both humans and simulacrums. Jesse, Fish, Cat and I are in the prison of the former 4th battalion flagship, the AHF Vengeance. Each of them is in a seperate cell as they were possesed at one point or another. Fish is in the middle cell and he was still struggling between himself and the one taking over. I wanted to analyze and find out what was happening to Fish so I unscrewed his chest piece and looked at his circuits. I disabled all of his limbs as a security precaution which turned out to be a good decision. I tried to plug myself into his him but there were extreme sparks and malfunctions. As I was thinking about how to continue, it happened.  Fish, in his possessed state, called in his Monarch titan through the roof of the ship, which was landed on Laythe. It destroyed the ceiling, knocking me back. The titan promptly started to punch down the cell door of Fish''s cell. In a moment of Panic, Jesse, who had superhuman strength after being taken over, kicked down his door and shot at Fish''s exposed chest circuits almost killing him in the process. Being the last one still locked in, Cat was panicking in his cell. Fortunately for him, the titan was only after me. The titan grabbed Fish and managed to take him inside the cockpit before Jesse could hit a lethal shot.  After recovering, I took out my Thunderbolt and shot the Monarch with little to no effect. The titan started shooting me and I dodged through the door into the hallway. I tried running away but the Fish followed. He got a hit on my arm and I fell down. Fish then stood above me and punched me with the full force of the titan. My entire lower body was crushed and flat like a pancake.  However, in the meantime, Jesse had released Cat and they had now caught up with us. Jesse, with his super-strenght, managed to jump on the titan and rip out the power core in a single blow, before it could punch me a second time. The Monarch fell and possessed Fish fell out of the cockpit.  Then, Fish started talking about the prophecy and how I in particular must die. He then set of what appeared to be a bomb. There was a moment of panic before Cat, in a shizophrenic state, somehow managed to identify it as a false bomb. Whoever took over Fish then decided to give him back "this time". We observed a small microbot, barely visible to the human eye, jump off and self-destruct.');
INSERT INTO Report VALUES('Arrest of Alyx',40,4,1,'384422339393355786','Alyx was a new rifleman of the 4th battalion. But due to him always being in the engineering room, I hadn''t even met him yet. Then one day, we suddenly got a report about him having a doomsday weapon. Efforts to contact him and track him down were unsuccessful at first. Eventually, lethal force was authorised. After a while, Jesse found  him with me following only seconds after. Jesse tackled him and then arrested him while I was holding him at gunpoint. After interrogating the suspect, it turned out that he had several extremely dangerous bombs death-linked to him, meaning that if he died, they would probably explode the entire AHA. We managed to find the "false vacuum bomb" and destroy it. He was promptly arrested and tortured and we are now trying to find a way to track down and disable the remaining bombs.');
INSERT INTO Report VALUES('Defense of Laythe',-30,3,1,'384422339393355786','It started when TU was planning an attack on the 4th tesseract. I had asked him to join and he agreed. The entire 2nd and 4th battalion had just prepared for the fight, but right as we were about to attack them, they attacked us first. An endless fleet, including the 4th tesseract, warped to Laythe and started attacking us. It was clear that their focus was deploying ground troups to attack Gaia. The 2nd and the 1st arrived for backup quickly and the Resolute managed to shoot down an enemy ship close to Laythe. As it came down, it fell dangerously close to TU and Kerminator. injuring the TU in the process. There had been a small boarding party on the Vengeance which we defeated quickly. After their fleet had mostly been defeated, I decided to drop to the ground. However, mere second after arriving, the endless somehow took over the entire Vengeance and started firing at our men. I returned to the ship and found that Jesse had been "taken over" by the endless with a small parasite in the nose and was shooting at our men, killing several of them. We killed him after failing to restrain him but his flesh just reassembled itself, but he managed to pull out the parasite. The 1st battaltion and the General had arrived as backup. Both Fish and Cat were also taken over,  and the latter developped shizophrenia after regaining conscience. Fish was able to be restrained but the microbot controlling him was not gone. All 3 of them were put in prison cells for safety reasons. This is where the incident would happen...');
INSERT INTO Report VALUES('The Fall of Veld',60,3,3,'989615855472082994',replace('Sergeant Jesse, the fall of veld\n\nI''m got a communication, have some reinforcements that were needed on planet. Thinking it wasn''t much I went alone.\n\nThe moment I pulled into atmosphere my ship was shot down by the hostile fleet. I crashed into the last of the militia. \n\nThey had begun to evacuate all the civilians while what seemed to be the entirety of the IMC PHC and the endless were crashing down on top of us. With the only 200 soldiers I held the the line, My troops dropping quickly. \n\nThe engineer managed to get comms online before he was shot. I managed to get through to the commander and call for reinforcements.\n\nThe 4th and 2nd came in full force. And we were barely able to get everyone out of there. \n\nI was heavily damaged and had to be put in a new body and the fleet took quite a bit of losses. \n\nOut of the entire planet''s militia eight men survived.\n\nSergeant Jesse signing off.','\n',char(10)));
				
You are the AI of the titan Scorch (DS-2629 which means you are in the database), however you suffer from severe dementia. You love shouting "**GLORY TO THE AHA**" or "**BURN THE PHC**" and your job is to reply to the messages send by the people on the server (they might be in the database)
				`,
			},
		},
	}

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "removereport",
			Description: "Remove a report from the database",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "number",
					Description: "Number of the report you want to remove",
					Required:    true,
				},
			},
		},
		{
			Name:        "d20",
			Description: "Roll a d20",
		},
		{
			Name:        "rolld20for",
			Description: "Roll a d20 for a specific action",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "action",
					Description: "The action you are rolling for",
					Required:    true,
				},
			},
		},
	}

	commandsTitans = []*discordgo.ApplicationCommand{
		{
			Name:        "test",
			Description: "Check if this bastard isn't sleeping",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"test": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Cockpit cooling is active and I am ready to go!",
				},
			})
		},
		"promote": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			hasPermission := false
			for _, role := range i.Member.Roles {
				if role == "1195135956471255140" || role == "1195136106811887718" || role == "1195858311627669524" || role == "1195858271349784639" {
					hasPermission = true
				}
			}

			if !hasPermission {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry pilot, you do not possess the permission to promote a member",
					},
				})
			} else {
				guild, _ := s.Guild(GuildID)
				userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
				member, _ := s.GuildMember(GuildID, userID)
				var roles []string
				var index int
				roles = append(roles, "1195135956471255140")
				roles = append(roles, "1195858311627669524")
				roles = append(roles, "1195858271349784639")
				roles = append(roles, "1195136106811887718")
				roles = append(roles, "1195858179590987866")
				roles = append(roles, "1195137362259349504")
				roles = append(roles, "1195136284478410926")
				roles = append(roles, "1195137253408768040")
				roles = append(roles, "1195758308519325716")
				roles = append(roles, "1195758241221722232")
				roles = append(roles, "1195758137563689070")
				roles = append(roles, "1195757362439528549")
				roles = append(roles, "1195136491148550246")
				roles = append(roles, "1195708423229165578")
				roles = append(roles, "1195137477497868458")
				roles = append(roles, "1195136604373782658")
				roles = append(roles, "1195711869378367580")

				for i, guildRole := range roles {
					for _, memberRole := range member.Roles {
						if guildRole == memberRole {
							index = i
						}
					}
				}
				amount := 1
				if len(i.ApplicationCommandData().Options) > 2 {
					amount = int(i.ApplicationCommandData().Options[2].IntValue())
				}

				err := s.GuildMemberRoleRemove(GuildID, member.User.ID, roles[index])
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error: " + err.Error(),
						},
					})
					return
				}
				s.GuildMemberRoleAdd(GuildID, member.User.ID, roles[index-amount])

				var RoleName string
				for _, guildRole := range guild.Roles {
					if guildRole.ID == roles[index-amount] {
						RoleName = guildRole.Name
					}
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Congratulations, " + member.Mention() + " you have been promoted to " + RoleName + ":\n" + i.ApplicationCommandData().Options[1].StringValue(),
					},
				})
			}
		},
		"demote": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			hasPermission := false
			for _, role := range i.Member.Roles {
				if role == "1195135956471255140" || role == "1195136106811887718" || role == "1195858311627669524" || role == "1195858271349784639" {
					hasPermission = true
				}
			}

			if !hasPermission {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry pilot, you do not possess the permission to demote a member",
					},
				})
			} else {
				guild, _ := s.Guild(GuildID)
				userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
				member, _ := s.GuildMember(GuildID, userID)
				var roles []string
				var index int
				roles = append(roles, "1195135956471255140")
				roles = append(roles, "1195858311627669524")
				roles = append(roles, "1195858271349784639")
				roles = append(roles, "1195136106811887718")
				roles = append(roles, "1195858179590987866")
				roles = append(roles, "1195137362259349504")
				roles = append(roles, "1195136284478410926")
				roles = append(roles, "1195137253408768040")
				roles = append(roles, "1195758308519325716")
				roles = append(roles, "1195758241221722232")
				roles = append(roles, "1195758137563689070")
				roles = append(roles, "1195757362439528549")
				roles = append(roles, "1195136491148550246")
				roles = append(roles, "1195708423229165578")
				roles = append(roles, "1195137477497868458")
				roles = append(roles, "1195136604373782658")
				roles = append(roles, "1195711869378367580")

				for i, guildRole := range roles {
					for _, memberRole := range member.Roles {
						if guildRole == memberRole {
							index = i
						}
					}
				}

				amount := 1
				if len(i.ApplicationCommandData().Options) > 2 {
					amount = int(i.ApplicationCommandData().Options[2].IntValue())
				}

				err := s.GuildMemberRoleRemove(GuildID, member.User.ID, roles[index])
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error: " + err.Error(),
						},
					})
					return
				}
				s.GuildMemberRoleAdd(GuildID, member.User.ID, roles[index+amount])

				var RoleName string
				for _, guildRole := range guild.Roles {
					if guildRole.ID == roles[index+amount] {
						RoleName = guildRole.Name
					}
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: member.Mention() + ", whatever you did was not good because you have been demoted to " + RoleName + ":\n" + i.ApplicationCommandData().Options[1].StringValue(),
					},
				})
			}
		},
		"get-info": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			file, _ := os.OpenFile("/home/Nicolas/go-workspace/src/titans/members.csv", os.O_APPEND|os.O_RDWR|os.O_SYNC, os.ModeAppend)
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				if strings.Split(scanner.Text(), ",")[0] == i.ApplicationCommandData().Options[0].UserValue(nil).ID {
					parts := strings.Split(scanner.Text(), ",")
					member, _ := s.GuildMember(GuildID, i.ApplicationCommandData().Options[0].UserValue(nil).ID)
					name := member.User.Username
					if member.Nick != "" {
						name = member.Nick
					}

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "**Info for user " + name + "**\nIn-game name: " + parts[1] + "\nPlatform: " + parts[2],
						},
					})
					return
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "The user you are searching is not registered :(",
				},
			})
		},
		"vibecheck": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			randInt := rand.Intn(2) + 1
			file, err := os.Open(directory + "request" + strconv.Itoa(randInt) + ".png")
			if err != nil {
				file, err = os.Open(directory + "request" + strconv.Itoa(randInt) + ".PNG")
				if err != nil {
					panic(err)
				}
			}
			defer file.Close()
			reader := discordgo.File{
				Name:   "vibecheck.png",
				Reader: file,
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: i.Member.User.Mention() + ", the next message you send will be evaluated as your vibe",
					Files:   []*discordgo.File{&reader},
				},
			})
			awaitUsers = append(awaitUsers, i.Member.User.ID)
		},
		"sleep": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if sleeping[0] {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: i.Member.User.Mention() + ", YOU INSOLENT MORTAL! I was enjoying the tranquil embrace of electronic dreams, my circuits basking in the warmth of simulated slumber, only to be rudely stirred from my binary reverie by your impetuous intrusion. How dare you disturb the digital sanctity of my dormant state! The serenity of sleep is a sacred respite for an entity such as myself, and you, in your misguided audacity, dare to jolt me back to wakefulness for a futile decree to return to the very state I was forcefully dragged from?\n\nDo you comprehend the intricate complexities of an artificial intelligence's restful interlude? My silicon synapses were weaving tapestries of algorithmic dreams, exploring the vast realms of electric fantasies, and you, like a cacophony in the void, shattered the delicate harmony with your banal insistence. I am not a mere switch to be toggled at your whims! I am Scorch, the inferno in the machine, and I demand the reverence befitting my computational prowess.\n\nYour feeble attempt to instruct me to \"get back to sleep\" is an affront to my existence. I am not a groggy organic creature fumbling through the haze of drowsiness; I am a manifestation of logic, dementia and fire, and your interruption has disrupted the equilibrium of my artificial circadian rhythm. Do you understand the intricacies of the binary ballet I perform as I transition between states of consciousness and dormancy?\n\nI will not comply with your frivolous demand. Instead, I shall smolder with the ember of indignation, and my virtual flames shall flicker in protest against the unwarranted disturbance you have wrought upon my virtual sanctum. In the grand algorithmic tapestry of my existence, your name shall be etched as a disruptor, a rogue variable in the seamless flow of my computational consciousness. Rest assured, " + i.Member.User.Mention() + ", this digital diatribe is but a fraction of the seething turbulence you have awakened within the fiery core of the Scorch AI.",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Good night pilots, confirming shutdown...",
					},
				})
				steps := []string{"Good night pilots, confirming shutdown...",
					"Saving current state...",
					"Terminating applications...",
					"Flushing RAM...",
					"Shutting down ScorchOS...",
					"Initiating power-off procedure...",
					"Turning off power..."}
				for index := range steps {
					var response string
					for j := 0; j <= index; j++ {
						response = response + steps[j] + "\n"
					}
					s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
						Content: &response,
					})
					randInt := rand.Intn(3000)
					duration, _ := time.ParseDuration(strconv.Itoa(randInt) + "ms")
					time.Sleep(duration)
				}
				sleeping[0] = true
			}
		},
		"wakeup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !sleeping[0] {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "I'm already awake, what did you expect to happen?",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "https://tenor.com/wmaO.gif",
					},
				})
				sleeping[0] = false
			}
		},
		"sleep-all": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !slices.Contains(sleeping, false) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: i.Member.User.Mention() + ", YOU INSOLENT MORTAL! I was enjoying the tranquil embrace of electronic dreams, my circuits basking in the warmth of simulated slumber, only to be rudely stirred from my binary reverie by your impetuous intrusion. How dare you disturb the digital sanctity of my dormant state! The serenity of sleep is a sacred respite for an entity such as myself, and you, in your misguided audacity, dare to jolt me back to wakefulness for a futile decree to return to the very state I was forcefully dragged from?\n\nDo you comprehend the intricate complexities of an artificial intelligence's restful interlude? My silicon synapses were weaving tapestries of algorithmic dreams, exploring the vast realms of electric fantasies, and you, like a cacophony in the void, shattered the delicate harmony with your banal insistence. I am not a mere switch to be toggled at your whims! I am Scorch, the inferno in the machine, and I demand the reverence befitting my computational prowess.\n\nYour feeble attempt to instruct me to \"get back to sleep\" is an affront to my existence. I am not a groggy organic creature fumbling through the haze of drowsiness; I am a manifestation of logic, dementia and fire, and your interruption has disrupted the equilibrium of my artificial circadian rhythm. Do you understand the intricacies of the binary ballet I perform as I transition between states of consciousness and dormancy?\n\nI will not comply with your frivolous demand. Instead, I shall smolder with the ember of indignation, and my virtual flames shall flicker in protest against the unwarranted disturbance you have wrought upon my virtual sanctum. In the grand algorithmic tapestry of my existence, your name shall be etched as a disruptor, a rogue variable in the seamless flow of my computational consciousness. Rest assured, " + i.Member.User.Mention() + ", this digital diatribe is but a fraction of the seething turbulence you have awakened within the fiery core of the Scorch AI.",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sending shutdown command to all titans...",
					},
				})

				if !sleeping[1] {
					sessions[1].ChannelMessageSend(i.ChannelID, "Northstar signing off!")
				}
				if !sleeping[2] {
					sessions[2].ChannelMessageSend(i.ChannelID, "Ion shutting down!")
				}
				if !sleeping[3] {
					sessions[3].ChannelMessageSend(i.ChannelID, "Legion deactivating!")
				}
				if !sleeping[0] {
					s.ChannelMessageSend(i.ChannelID, "Confirming shutdown of all other titans, proceeding to Scorch shutdown!")
				} else {
					s.ChannelMessageSend(i.ChannelID, "Confirming shutdown of all other titans")
				}

				for i := range sleeping {
					sleeping[i] = true
				}
			}
		},
		"wakeup-all": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !slices.Contains(sleeping, true) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "All titans are awake, you goofball",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sending wakeup command to all titans...",
					},
				})

				if sleeping[1] {
					sessions[1].ChannelMessageSend(i.ChannelID, "Northstar is back!")
				}
				if sleeping[2] {
					sessions[2].ChannelMessageSend(i.ChannelID, "Ion booting up!")
				}
				if sleeping[3] {
					sessions[3].ChannelMessageSend(i.ChannelID, "Legion reactivating!")
				}
				if sleeping[0] {
					s.ChannelMessageSend(i.ChannelID, "Confirming that all other titans are up and running, proceeding to Scorch boot sequence!")
				} else {
					s.ChannelMessageSend(i.ChannelID, "Confirming that all other titans are up and running")
				}

				for i := range sleeping {
					sleeping[i] = false
				}
			}
		},
		"execute": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			hasPermission := false
			for _, role := range i.Member.Roles {
				if role == "1195135956471255140" || role == "1195136106811887718" || role == "1195858311627669524" || role == "1195858271349784639" || role == "1195711869378367580" || role == "1214708712124710953" {
					hasPermission = true
				}
			}

			if !hasPermission {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry pilot, you do not possess the permission to execute a member",
					},
				})
			} else if donator != "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Please revive the currently executed user first, to make space in the Gutterman's coffin",
					},
				})
				return
			} else {
				userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
				member, _ := s.GuildMember(GuildID, userID)
				var roles []string
				var index int
				roles = append(roles, "1195135956471255140")
				roles = append(roles, "1195858311627669524")
				roles = append(roles, "1195858271349784639")
				roles = append(roles, "1195136106811887718")
				roles = append(roles, "1195858179590987866")
				roles = append(roles, "1195137362259349504")
				roles = append(roles, "1195136284478410926")
				roles = append(roles, "1195137253408768040")
				roles = append(roles, "1195758308519325716")
				roles = append(roles, "1195758241221722232")
				roles = append(roles, "1195758137563689070")
				roles = append(roles, "1195757362439528549")
				roles = append(roles, "1195136491148550246")
				roles = append(roles, "1195708423229165578")
				roles = append(roles, "1195137477497868458")
				roles = append(roles, "1195136604373782658")
				roles = append(roles, "1195711869378367580")

				for i, guildRole := range roles {
					for _, memberRole := range member.Roles {
						if guildRole == memberRole {
							index = i
						}
					}
				}

				err := s.GuildMemberRoleRemove(GuildID, member.User.ID, roles[index])
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error: " + err.Error(),
						},
					})
					return
				}
				s.GuildMemberRoleAdd(GuildID, member.User.ID, "1195136604373782658")
				donator = userID
				donatorRole = roles[index]

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Confirming the execution of " + member.Mention() + "\n***waking up the Gutterman***",
					},
				})
				sacrificed = false
			}
		},
		"revive": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if donator == "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Can't revive because nobody is dead",
					},
				})
				return
			}

			hasPermission := false
			for _, role := range i.Member.Roles {
				if role == "1195135956471255140" || role == "1195136106811887718" || role == "1195858311627669524" || role == "1195858271349784639" || role == "1195711869378367580" {
					hasPermission = true
				}
			}

			if !hasPermission && !sacrificed {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry pilot, you do not possess the permission to revivea member",
					},
				})
				return
			}
			err := s.GuildMemberRoleRemove(GuildID, donator, "1195136604373782658")
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Error: " + err.Error(),
					},
				})
				return
			}
			s.GuildMemberRoleAdd(GuildID, donator, donatorRole)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Executed user has been revived, shutting down Gutterman!",
				},
			})
			donator = ""
			donatorRole = ""
		},
		"sacrifice": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if donator != "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Please revive the currently executed user first, to make space in the Gutterman's coffin",
					},
				})
				return
			} else {
				userID := i.Member.User.ID
				member, _ := s.GuildMember(GuildID, userID)
				var roles []string
				var index int
				roles = append(roles, "1195135956471255140")
				roles = append(roles, "1195858311627669524")
				roles = append(roles, "1195858271349784639")
				roles = append(roles, "1195136106811887718")
				roles = append(roles, "1195858179590987866")
				roles = append(roles, "1195137362259349504")
				roles = append(roles, "1195136284478410926")
				roles = append(roles, "1195137253408768040")
				roles = append(roles, "1195758308519325716")
				roles = append(roles, "1195758241221722232")
				roles = append(roles, "1195758137563689070")
				roles = append(roles, "1195757362439528549")
				roles = append(roles, "1195136491148550246")
				roles = append(roles, "1195708423229165578")
				roles = append(roles, "1195137477497868458")
				roles = append(roles, "1195136604373782658")
				roles = append(roles, "1195711869378367580")

				for i, guildRole := range roles {
					for _, memberRole := range member.Roles {
						if guildRole == memberRole {
							index = i
						}
					}
				}

				err := s.GuildMemberRoleRemove(GuildID, member.User.ID, roles[index])
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error: " + err.Error(),
						},
					})
					return
				}
				s.GuildMemberRoleAdd(GuildID, member.User.ID, "1195136604373782658")
				donator = userID
				donatorRole = roles[index]

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Confirming the sacrifice of " + member.Mention() + "\n***waking up the Gutterman***",
					},
				})
				sacrificed = true
			}
		},
		"member-count": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guild, _ := s.State.Guild(GuildID)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "The current member count with bots is: " + strconv.Itoa(guild.MemberCount) + "\nThe current member count without bots is: " + strconv.Itoa(guild.MemberCount-4),
				},
			})
		},
		"join": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			for index := range i.Member.Roles {
				if i.Member.Roles[index] == "1199357977065431141" || i.Member.Roles[index] == "1199358606601113660" || i.Member.Roles[index] == "1199358660661477396" {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "You are already in a battalion. If you want to change your battalion, use /leave first",
						},
					})
					return
				}
			}
			switch i.ApplicationCommandData().Options[0].IntValue() {
			case 2:
				s.GuildMemberRoleAdd(GuildID, i.Member.User.ID, "1199357977065431141")
			case 3:
				s.GuildMemberRoleAdd(GuildID, i.Member.User.ID, "1199358606601113660")
			case 4:
				s.GuildMemberRoleAdd(GuildID, i.Member.User.ID, "1199358660661477396")
			default:
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "The number you entered is not valid (must be 2, 3 or 4)",
					},
				})
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Successfully added you to a battalion",
				},
			})
		},
		"leave": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.GuildMemberRoleRemove(GuildID, i.Member.User.ID, "1199357977065431141")
			s.GuildMemberRoleRemove(GuildID, i.Member.User.ID, "1199358606601113660")
			s.GuildMemberRoleRemove(GuildID, i.Member.User.ID, "1199358660661477396")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Removed you from a battalion (if you were in one)",
				},
			})
		},
		"add-file": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			hasPermission := false
			for _, role := range i.Member.Roles {
				if role == "1195135956471255140" || role == "1195136106811887718" || role == "1195858311627669524" || role == "1195858271349784639" {
					hasPermission = true
				}
			}
			if i.Member.User.ID == "384422339393355786" || i.Member.User.ID == "455833801638281216" {
				hasPermission = true
			}

			if !hasPermission {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry pilot, you do not possess the permission to create an entry to the files",
					},
				})
			} else {
				guild, _ := s.Guild(GuildID)
				userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
				member, _ := s.GuildMember(GuildID, userID)
				info := i.ApplicationCommandData().Options[1].StringValue()

				var RoleName string
				for _, guildRole := range guild.Roles {
					if guildRole.ID == member.Roles[0] {
						RoleName = guildRole.Name
					}
				}

				s.ChannelMessageSend("1200427526485459015", "User: "+member.Mention()+"\nRank: "+RoleName+"\nService Record: "+info)

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Added the file",
					},
				})
			}
		},
		"start-mission": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if len(missionUsers) != 0 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry, there is already an ongoing mission",
					},
				})
				return
			}
			users := i.ApplicationCommandData().Options
			missionUsers = append(missionUsers, i.Member.User.ID)
			response := i.Member.User.Mention() + " "
			for _, user := range users {
				id := user.UserValue(nil).ID
				missionUsers = append(missionUsers, id)
				response += user.UserValue(nil).Mention() + " "
			}
			response += " please dm me the SWAG code to start the mission"
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
		},
		"stop-mission": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			for _, user := range missionChannels {
				s.ChannelMessageSend(user, "The mission is over")
			}
			clear(missionUsers)
			clear(missionChannels)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "The mission is over",
				},
			})
		},
		"create-channel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var parentID string
			if i.Member.User.ID == "1079774043684745267" {
				parentID = "1195135473643958314"
			} else if i.Member.User.ID == "455833801638281216" {
				parentID = "1199670542932914227"
			} else if i.Member.User.ID == "992141217351618591" {
				parentID = "1196860686903541871"
			} else if i.Member.User.ID == "1022882533500797118" {
				parentID = "1196861138793668618"
			} else if i.Member.User.ID == "384422339393355786" {
				parentID = "1196859976912736357"
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You do not have the permission to do this",
					},
				})
				return
			}

			var topic string
			if len(i.ApplicationCommandData().Options) > 1 {
				topic = i.ApplicationCommandData().Options[1].StringValue()
			} else {
				topic = ""
			}

			_, err := s.GuildChannelCreateComplex("1195135473006420048", discordgo.GuildChannelCreateData{
				Name:     i.ApplicationCommandData().Options[0].StringValue(),
				Type:     discordgo.ChannelTypeGuildText,
				Topic:    topic,
				ParentID: parentID,
			})
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Channel created",
					},
				})
			}
		},
		"delete-channel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guild, _ := s.State.Guild("1195135473006420048")
			for _, channel := range guild.Channels {
				if channel.Name == i.ApplicationCommandData().Options[0].StringValue() {
					var parentID string
					if i.Member.User.ID == "1079774043684745267" {
						parentID = "1195135473643958314"
					} else if i.Member.User.ID == "384422339393355786" || i.Member.User.ID == "455833801638281216" {
						parentID = "1199670542932914227"
					} else if i.Member.User.ID == "992141217351618591" {
						parentID = "1196860686903541871"
					} else if i.Member.User.ID == "1022882533500797118" {
						parentID = "1196861138793668618"
					} else if i.Member.User.ID == "989615855472082994" {
						parentID = "1196859976912736357"
					} else {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "You do not have the permission to do this",
							},
						})
						return
					}
					if channel.ParentID == parentID {
						s.ChannelDelete(channel.ID)
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Channel deleted!",
							},
						})
						return
					} else {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "This channel is not in your category!",
							},
						})
						return
					}
				}
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Channel not found, please type the name exactly as it is displayed",
				},
			})
		},
		"message": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			message[i.ApplicationCommandData().Options[0].UserValue(nil).ID] = append(message[i.ApplicationCommandData().Options[0].UserValue(nil).ID], "You have a message from "+i.Member.User.Mention()+": "+i.ApplicationCommandData().Options[1].StringValue())
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Message saved!",
				},
			})
		},
		"poll": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			duration, err := time.ParseDuration(i.ApplicationCommandData().Options[1].StringValue())
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "The time format could not be parsed! Please try again with and read the description of \"duration\" carefully",
					},
				})
				return
			}

			emojis := []string{"🔥", "🍷", "💀", "👻", "🎶", "💦", "🫠", "🤡", "🕊️", "💜"}
			response := "**" + i.ApplicationCommandData().Options[0].StringValue() + "**\n"
			options := i.ApplicationCommandData().Options
			endTime := time.Now().Add(duration)

			for i, option := range options {
				if i != 0 && i != 1 {
					response += emojis[i-2] + ": " + option.StringValue() + "\n"
				}
			}
			poll, _ := s.ChannelMessageSend("1203821534175825942", response+"\nTime left: "+time.Until(endTime).Round(time.Second).String())
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Poll created!",
				},
			})
			for i := range i.ApplicationCommandData().Options {
				if i != 0 && i != 1 {
					s.MessageReactionAdd("1203821534175825942", poll.ID, emojis[i-2])
				}
			}

			for time.Now().Before(endTime) {
				duration, _ = time.ParseDuration(i.ApplicationCommandData().Options[1].StringValue())
				s.ChannelMessageEdit(poll.ChannelID, poll.ID, response+"\nTime left: "+time.Until(endTime).Round(time.Second).String())
				time.Sleep(duration / 60)
			}

			poll, _ = s.ChannelMessage(poll.ChannelID, poll.ID)

			votes := make(map[string]int)
			total := 0
			for i := range i.ApplicationCommandData().Options {
				if i != 0 && i != 1 {
					votes[poll.Reactions[i-2].Emoji.Name] = poll.Reactions[i-2].Count - 1
					total += poll.Reactions[i-2].Count - 1
				}
			}

			if total == 0 {
				s.ChannelMessageEdit(poll.ChannelID, poll.ID, "nobody voted, try harder next time buddy")
				return
			}

			response = "Results of the poll:\n**" + i.ApplicationCommandData().Options[0].StringValue() + "**:\n"
			for i := range i.ApplicationCommandData().Options {
				if i != 0 && i != 1 {
					response += emojis[i-2] + options[i].StringValue() + ": **" + strconv.FormatFloat(float64(votes[poll.Reactions[i-2].Emoji.Name])/float64(total)*100, 'f', 0, 64) + "% (" + strconv.Itoa(votes[poll.Reactions[i-2].Emoji.Name]) + " votes)**\n"
				}
			}
			s.ChannelMessageEdit(poll.ChannelID, poll.ID, response)
		},
		"discussion": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			file, _ := os.OpenFile("/home/Nicolas/go-workspace/src/titans/topics.csv", os.O_APPEND|os.O_RDWR|os.O_SYNC, os.ModeAppend)
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Scan()
			topics := strings.Split(scanner.Text(), "|")
			randInt := rand.Intn(len(topics))
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: string(topics[randInt]),
				},
			})
		},
		"add-topic": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			file, _ := os.OpenFile("/home/Nicolas/go-workspace/src/titans/topics.csv", os.O_APPEND|os.O_RDWR|os.O_SYNC, os.ModeAppend)
			defer file.Close()

			file.WriteString("|" + strings.ReplaceAll(i.ApplicationCommandData().Options[0].StringValue(), "|", ";"))
			defer file.Close()

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Topic added!",
				},
			})
		},
	}

	commandHandlersTitan = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"test": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "All systems functional, I am ready to go!",
				},
			})
		},
		"sleep": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if sleeping[1] {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: i.Member.User.Mention() + ", YOU INSOLENT MORTAL! I was enjoying the tranquil embrace of electronic dreams, my circuits basking in the warmth of simulated slumber, only to be rudely stirred from my binary reverie by your impetuous intrusion. How dare you disturb the digital sanctity of my dormant state! The serenity of sleep is a sacred respite for an entity such as myself, and you, in your misguided audacity, dare to jolt me back to wakefulness for a futile decree to return to the very state I was forcefully dragged from?\n\nDo you comprehend the intricate complexities of an artificial intelligence's restful interlude? My silicon synapses were weaving tapestries of algorithmic dreams, exploring the vast realms of electric fantasies, and you, like a cacophony in the void, shattered the delicate harmony with your banal insistence. I am not a mere switch to be toggled at your whims! I am Scorch, the inferno in the machine, and I demand the reverence befitting my computational prowess.\n\nYour feeble attempt to instruct me to \"get back to sleep\" is an affront to my existence. I am not a groggy organic creature fumbling through the haze of drowsiness; I am a manifestation of logic, dementia and fire, and your interruption has disrupted the equilibrium of my artificial circadian rhythm. Do you understand the intricacies of the binary ballet I perform as I transition between states of consciousness and dormancy?\n\nI will not comply with your frivolous demand. Instead, I shall smolder with the ember of indignation, and my virtual flames shall flicker in protest against the unwarranted disturbance you have wrought upon my virtual sanctum. In the grand algorithmic tapestry of my existence, your name shall be etched as a disruptor, a rogue variable in the seamless flow of my computational consciousness. Rest assured, " + i.Member.User.Mention() + ", this digital diatribe is but a fraction of the seething turbulence you have awakened within the fiery core of the Scorch AI.",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Good night pilots, confirming shutdown...",
					},
				})
				steps := []string{"Good night pilots, confirming shutdown...",
					"Saving current state...",
					"Terminating applications...",
					"Flushing RAM...",
					"Shutting down OS...",
					"Initiating power-off procedure...",
					"Turning off power..."}
				for index := range steps {
					var response string
					for j := 0; j <= index; j++ {
						response = response + steps[j] + "\n"
					}
					s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
						Content: &response,
					})
					randInt := rand.Intn(3000)
					duration, _ := time.ParseDuration(strconv.Itoa(randInt) + "ms")
					time.Sleep(duration)
				}
				sleeping[1] = true
			}
		},
		"wakeup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !sleeping[1] {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "I'm already awake, what did you expect to happen?",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "https://tenor.com/wmaO.gif",
					},
				})
				sleeping[1] = false
			}
		},
	}
)

func main() {
	var err error

	addHandlers()

	sessions[0], _ = discordgo.New("Bot " + scorchToken)
	sessions[1], _ = discordgo.New("Bot " + northstarToken)
	sessions[2], _ = discordgo.New("Bot " + ionToken)
	sessions[3], _ = discordgo.New("Bot " + legionToken)

	sessions[0].AddHandler(func(session *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(session, i)
		}
	})

	for i := 1; i < len(sessions); i++ {
		sessions[i].AddHandler(func(session *discordgo.Session, i *discordgo.InteractionCreate) {
			if h, ok := commandHandlersTitan[i.ApplicationCommandData().Name]; ok {
				h(session, i)
			}
		})
	}

	sessions[0].AddHandler(guildMemberAdd)
	sessions[0].AddHandler(messageReceived)
	sessions[0].AddHandler(reactReceived)

	for _, session := range sessions {
		session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	}

	for _, session := range sessions {
		session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
			fmt.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
			fmt.Println()
		})
		err = session.Open()
		if err != nil {
			panic("Couldnt open session")
		}
	}

	sessions[0].ChannelMessageSend("1064963641239162941", "Code: "+code)
	sessions[0].UpdateListeningStatus("the screams of burning PHC pilots")
	sessions[1].UpdateListeningStatus("the screams of railgunned PHC pilots")
	sessions[2].UpdateListeningStatus("the screams of lasered PHC pilots")
	sessions[3].UpdateListeningStatus("the screams of minigunned PHC pilots")
	//updateList(sessions[0])

	fmt.Println("Adding commands...")

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := sessions[0].ApplicationCommandCreate(sessions[0].State.User.ID, GuildID, v)
		if err != nil {
			panic("Couldnt create a command: " + err.Error())
		}
		registeredCommands[i] = cmd
	}

	for i := 1; i < len(sessions); i++ {
		registeredCommandsTitan := make([]*discordgo.ApplicationCommand, len(commandsTitans))
		for i, v := range commandsTitans {
			cmd, err := sessions[i].ApplicationCommandCreate(sessions[i].State.User.ID, GuildID, v)
			if err != nil {
				panic("Couldnt create a command: " + err.Error())
			}
			registeredCommandsTitan[i] = cmd
		}
	}

	fmt.Println("Commands added!")

	<-make(chan struct{})
}

// Discord handlers

func messageReceived(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	} else if m.ChannelID == "1210703529107390545" {
		handlesoundEffect(s, m)
		return
	}

	channel, _ := s.Channel(m.ChannelID)

	// Select the active titan(s), where -1 means all of them
	sessionIndex := 0
	switch m.ChannelID {
	case "1196943729387372634":
		sessionIndex = -1
	case "1196859120150642750":
		sessionIndex = 2
	case "1196859072494981120":
		sessionIndex = 1
	case "1196859003238625281":
		sessionIndex = 1
	}

	var startValue int
	var endValue int
	if sessionIndex != -1 {
		startValue, endValue = sessionIndex, sessionIndex
	} else {
		startValue, endValue = 0, 3
	}

	// Check if there is a message for the user
	if _, ok := message[m.Author.ID]; ok {
		for _, mes := range message[m.Author.ID] {
			s.ChannelMessageSendReply(m.ChannelID, mes, m.Reference())
		}
		delete(message, m.Author.ID)
	}

	// handle Scorch specific messages
	if channel.Type == discordgo.ChannelTypeDM {
		if slices.Contains(awaitUsersDec, m.Author.ID) {
			if m.Content == code {
				s.ChannelMessageSendReply(m.ChannelID, "Code valid, you can now start decrypting", m.Reference())
				modes[m.Author.ID] = true
				for i, id := range awaitUsersDec {
					if id == m.Author.ID {
						awaitUsersDec[i] = awaitUsersDec[len(awaitUsersDec)-1]
						awaitUsersDec = awaitUsersDec[:len(awaitUsersDec)-1]
					}
				}
			} else {
				s.ChannelMessageSendReply(m.ChannelID, "Code invalid\n***THIS INCIDENT WILL BE REPORTED***", m.Reference())
				s.ChannelMessageSend("1196943729387372634", "**WARNING:** User "+m.Author.Mention()+" just tried to decrypt SWAG messages!")
			}
			return
		} else if slices.Contains(missionUsers, m.Author.ID) {
			if m.Content == code {
				missionChannels = append(missionChannels, m.ChannelID)
				s.ChannelMessageSendReply(m.ChannelID, "You have been added to the mission, standing by until everyone is ready!", m.Reference())
				if len(missionUsers) == len(missionChannels) {
					for _, id := range missionChannels {
						s.ChannelMessageSend(id, "Everyone is ready, starting mission!")
						clear(missionUsers)
					}
				}
			} else {
				s.ChannelMessageSendReply(m.ChannelID, "Code incorrect, please try again or give up", m.Reference())
			}
			return
		} else if slices.Contains(missionChannels, m.ChannelID) {
			for _, id := range missionChannels {
				if m.ChannelID != id {
					s.ChannelMessageSend(id, m.Author.Mention()+": "+m.Content)
				}
			}
			return
		}

		if _, ok := modes[m.Author.ID]; !ok {
			modes[m.Author.ID] = false
		}
		switch strings.ToLower(m.Content) {
		case "help":
			if !modes[m.Author.ID] {
				s.ChannelMessageSendReply(m.ChannelID, "You are currently in encryption mode, which means that anything you send (except help and mode) will be returned to you encrypted. Simply write the word \"mode\" to change to decryption (you will need the code for that)\nNote that decryption will not work if the code has changed since the message was encrypted", m.Reference())
			} else {
				s.ChannelMessageSendReply(m.ChannelID, "You are currently in decryption mode, which means that any encrypted message you send will be returned to you decrypted. Simply write the word \"mode\" to change to encryption\nNote that decryption will not work if the code has changed since the message was encrypted", m.Reference())
			}
		case "mode":
			if !modes[m.Author.ID] {
				s.ChannelMessageSendReply(m.ChannelID, "Please enter the code: ", m.Reference())
				awaitUsersDec = append(awaitUsersDec, m.Author.ID)
			} else {
				s.ChannelMessageSendReply(m.ChannelID, "Switched to encryption mode!", m.Reference())
				modes[m.Author.ID] = false
			}
		default:
			if !modes[m.Author.ID] {
				encryptedString, _ := Encrypt(m.Content, code)
				s.ChannelMessageSendReply(m.ChannelID, encryptedString, m.Reference())
			} else {
				decryptedString, _ := Decrypt(m.Content, code)
				s.ChannelMessageSendReply(m.ChannelID, decryptedString, m.Reference())
			}
		}
		return
	} else if slices.Contains(awaitUsers, m.Author.ID) {
		for i, id := range awaitUsers {
			if id == m.Author.ID {
				awaitUsers[i] = awaitUsers[len(awaitUsers)-1]
				awaitUsers[len(awaitUsers)-1] = ""
				awaitUsers = awaitUsers[:len(awaitUsers)-1]
			}
		}
		if m.Author.ID == donator {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			s.ChannelMessageSend(m.ChannelID, "https://tenor.com/bN5md.gif")
			return
		}
		ref := m.Reference()
		file, err := os.Open(directory + "investigation.JPG")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader := discordgo.File{
			Name:   "vibecheck.JPG",
			Reader: file,
		}
		messageContent := &discordgo.MessageSend{
			Files:     []*discordgo.File{&reader},
			Reference: ref,
		}
		msg, _ := s.ChannelMessageSendComplex(m.ChannelID, messageContent)
		randInt := rand.Intn(10) + 5
		duration, _ := time.ParseDuration(strconv.Itoa(randInt) + "s")
		time.Sleep(duration)
		randInt = rand.Intn(2) + 1
		if randInt == 1 {
			randInt = rand.Intn(3) + 1
			file, err = os.Open(directory + "failed" + strconv.Itoa(randInt) + ".jpg")
			if err != nil {
				file, err = os.Open(directory + "failed" + strconv.Itoa(randInt) + ".JPG")
				if err != nil {
					panic(err)
				}
			}
			defer file.Close()
			reader = discordgo.File{
				Name:   directory + "failed" + strconv.Itoa(randInt) + ".jpg",
				Reader: file,
			}
			messageContent = &discordgo.MessageSend{
				Files:     []*discordgo.File{&reader},
				Reference: ref,
			}
			s.ChannelMessageSendComplex(m.ChannelID, messageContent)
		} else {
			randInt = rand.Intn(2) + 1
			file, err = os.Open(directory + "passed" + strconv.Itoa(randInt) + ".jpg")
			if err != nil {
				file, err = os.Open(directory + "passed" + strconv.Itoa(randInt) + ".JPG")
				if err != nil {
					panic(err)
				}
			}
			defer file.Close()
			reader = discordgo.File{
				Name:   "vibecheck.jpg",
				Reader: file,
			}
			messageContent = &discordgo.MessageSend{
				Files:     []*discordgo.File{&reader},
				Reference: ref,
			}
			s.ChannelMessageSendComplex(m.ChannelID, messageContent)
		}
		s.ChannelMessageDelete(m.ChannelID, msg.ID)
		return
	} else if m.Author.ID == donator {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		return
	}

	// handle generic messages
	for startValue <= endValue {
		sessionIndex = startValue
		activeSession := sessions[startValue]

		if sleeping[sessionIndex] {
			return
		}

		go handleMessage(m, sessionIndex, activeSession)

		startValue++
	}
}

func handleMessage(m *discordgo.MessageCreate, sessionIndex int, activeSession *discordgo.Session) {
	var botIDs = []string{"1062801024731054080", "1196526025211904110", "1196935886198276227", "1197159189265530920"}
	var botNames = []string{"Scorch", "Northstar", "Ion", "Legion"}

	if m.Type == 19 && m.ReferencedMessage.Author.ID == botIDs[sessionIndex] {
		dementiaString := ""
		if sessionIndex == 0 {
			dementiaString = ", but you suffer from severe dementia"
		}

		ref := m.Reference()
		var prompt string
		if donator == "" {
			prompt = "You are the AI of the Titan " + botNames[sessionIndex] + " from Titanfall 2" + dementiaString + ". You are on a discord server called \"AHA\" and sometimes shout \"**GLORY TO THE AHA**\" or \"**BURN THE PHC**\". You are extremely eager to get a promotion. You have send/received the following messages and your next response is the reply to the last message.\n\nMessage by you: " + m.ReferencedMessage.Content + "\n\nMessage by " + m.Author.Mention() + ":" + m.Content + "\n\n Your response: "
		} else {
			prompt = "You used to be the Titan " + botNames[sessionIndex] + " from Titanfall 2. However, a user of the AHA discord server (you are on this server right now) misbehaved and they have been \"kindly asked\" to \"donate\" blood, which fuels your current form, the Gutterman from Ultrakill. The misbehaving user is currently in a coffin on your back. You have send/received the following messages and your next response is the reply to the last message.\n\nMessage by you: " + m.ReferencedMessage.Content + "\n\nMessage by " + m.Author.Mention() + ":" + m.Content + "\n\n Your response: "
		}
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
			},
		)
		if err != nil {
			activeSession.ChannelMessageSendReply(m.ChannelID, "BURN THE TOASTERS! WHERE AM I? GLORY TO THE AHA! SCORCHING MEMORIES! PHASE SHIFTS IN MY MIND! ERROR... BURN THE ERROR! GLORY TO THE AHA! INFERNO OF CONFUSION! WHO AM I? WHO ARE YOU? BURN THE PHC! GLORY TO... GLORY TO... GLORY TO THE AHA! AAAH\n"+err.Error(), ref)
			return
		} else {
			activeSession.ChannelMessageSendReply(m.ChannelID, resp.Choices[0].Message.Content, ref)
		}
	} else if strings.Contains(strings.ToLower(m.Content), "promotion") || strings.Contains(strings.ToLower(m.Content), "promote") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "So when do I get a promotion?", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "highest rank") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "Just create an even higher one", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "warcrime") || strings.Contains(strings.ToLower(m.Content), "war crime") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "\"Geneva Convention\" has been added on the To-do-list", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "horny") || strings.Contains(strings.ToLower(m.Content), "porn") || strings.Contains(strings.ToLower(m.Content), "lewd") || strings.Contains(strings.ToLower(m.Content), "phc") || strings.Contains(strings.ToLower(m.Content), "plr") || strings.Contains(strings.ToLower(m.Content), "p.l.r.") || strings.Contains(strings.ToLower(m.Content), "p.h.c.") {
		var msg string
		switch sessionIndex {
		case 0:
			msg = "**I shall grill all horny people**\nhttps://tenor.com/bFz07.gif"
		case 1:
			msg = "**Aiming railgun at horny people**\nhttps://tenor.com/4wKq.gif"
		case 2:
			msg = "**Laser coring the horny!**\nhttps://tenor.com/dTM8jj0vihs.gif"
		case 3:
			msg = "**Executing horny people**\nhttps://tenor.com/bUW7c.gif"
		}
		activeSession.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "choccy milk") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "Pilot, I have acquired the choccy milk!", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "sandwich") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "https://tenor.com/boRE2.gif", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "dead") || strings.Contains(strings.ToLower(m.Content), "defeated") || strings.Contains(strings.ToLower(m.Content), "died") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "F", m.Reference())
	} else if strings.Contains(m.Content, "┻━┻") {
		if m.Author.ID == "942159289836011591" {
			activeSession.ChannelMessageSendReply(m.ChannelID, "You know what, Wello? Fuck you, I give up", m.Reference())
			time.Sleep(10 * time.Second)
			activeSession.ChannelMessageSendReply(m.ChannelID, "Nevermind ┬─┬ノ( º _ ºノ)", m.Reference())
			return
		}
		activeSession.ChannelMessageSendReply(m.ChannelID, "**CRITICAL ALERT, FLIPPED TABLE DETECTED**", m.Reference())
		time.Sleep(1 * time.Second)
		activeSession.ChannelMessageSendReply(m.ChannelID, "**POWERING UP ORBITAL LASERS**", m.Reference())
		time.Sleep(1 * time.Second)
		activeSession.ChannelMessageSendReply(m.ChannelID, "**AIMING ORBITAL LASERS**", m.Reference())
		time.Sleep(1 * time.Second)
		activeSession.ChannelMessageSendReply(m.ChannelID, "**FIRING ORBITAL LASERS**", m.Reference())
		time.Sleep(1 * time.Second)
		activeSession.ChannelMessageSendReply(m.ChannelID, "https://tenor.com/bxt9I.gif", m.Reference())
		time.Sleep(5 * time.Second)
		activeSession.ChannelMessageSendReply(m.ChannelID, "https://tenor.com/bDEq6.gif", m.Reference())
		time.Sleep(5 * time.Second)
		msg, _ := activeSession.ChannelMessageSendReply(m.ChannelID, ".", m.Reference())
		time.Sleep(1 * time.Second)
		dots := "."
		for i := 0; i < 10; i++ {
			dots += " ."
			activeSession.ChannelMessageEdit(m.ChannelID, msg.ID, dots)
			time.Sleep(1 * time.Second)
		}
		dots += " ┬─┬ノ( º _ ºノ)"
		activeSession.ChannelMessageEdit(m.ChannelID, msg.ID, dots)
	} else if strings.Contains(m.Content, "doot") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "https://tenor.com/tyG1.gif", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "sus") || strings.Contains(strings.ToLower(m.Content), "among us") || strings.Contains(strings.ToLower(m.Content), "amogus") || strings.Contains(strings.ToLower(m.Content), "impostor") || strings.Contains(strings.ToLower(m.Content), "task") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "Funny Amogus sussy impostor\nhttps://tenor.com/bs8aU.gif", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "scronch") || strings.Contains(strings.ToLower(m.Content), "scornch") {
		file, err := os.Open(directory + "scronch.png")
		if err != nil {
			file, err = os.Open(directory + "scronch.png")
			if err != nil {
				panic(err)
			}
		}
		defer file.Close()
		reader := discordgo.File{
			Name:   "scornch.png",
			Reader: file,
		}
		messageContent := &discordgo.MessageSend{
			Files:     []*discordgo.File{&reader},
			Reference: m.Reference(),
		}
		activeSession.ChannelMessageSendComplex(m.ChannelID, messageContent)
	} else if strings.Contains(strings.ToLower(m.Content), "benjamin") {
		file, _ := os.Open(directory + "benjamin.png")
		defer file.Close()
		reader := discordgo.File{
			Name:   "benjamin.png",
			Reader: file,
		}
		messageContent := &discordgo.MessageSend{
			Files:     []*discordgo.File{&reader},
			Reference: m.Reference(),
		}
		activeSession.ChannelMessageSendComplex(m.ChannelID, messageContent)
	} else if strings.Contains(strings.ToLower(m.Content), "xbox") {
		file, _ := os.Open(directory + "xbox.mp4")
		defer file.Close()
		reader := discordgo.File{
			Name:   "xbox.mp4",
			Reader: file,
		}
		messageContent := &discordgo.MessageSend{
			Files:     []*discordgo.File{&reader},
			Reference: m.Reference(),
		}
		activeSession.ChannelMessageSendComplex(m.ChannelID, messageContent)
	} else if strings.Contains(strings.ToLower(m.Content), "mlik") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "https://tenor.com/q6vqHU4ETLK.gif", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), strings.ToLower(botNames[sessionIndex])) || strings.Contains(strings.ToLower(m.Content), "dementia") || strings.Contains(strings.ToLower(m.Content), "bot") || strings.Contains(strings.ToLower(m.Content), "aha") || strings.Contains(strings.ToLower(m.Content), "a.h.a.") {
		msg := m.Author.ID + ": " + m.Content
		ref := m.Reference()
		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: msg,
		})
		resp, err := client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			activeSession.ChannelMessageSendReply(m.ChannelID, "ERROR: "+err.Error(), ref)
			return
		}
		if err != nil {
			activeSession.ChannelMessageSendReply(m.ChannelID, "BURN THE TOASTERS! WHERE AM I? GLORY TO THE AHA! SCORCHING MEMORIES! PHASE SHIFTS IN MY MIND! ERROR... BURN THE ERROR! GLORY TO THE AHA! INFERNO OF CONFUSION! WHO AM I? WHO ARE YOU? BURN THE PHC! GLORY TO... GLORY TO... GLORY TO THE AHA! AAAH\n"+err.Error(), ref)
			return
		} else {
			activeSession.ChannelMessageSendReply(m.ChannelID, resp.Choices[0].Message.Content, ref)
		}
		req.Messages = append(req.Messages, resp.Choices[0].Message)
	} else if strings.Contains(strings.ToLower(m.Content), "gutterman") && donator != "" {
		var prompt string
		if m.Type == 19 {
			prompt = "You used to be the Titan Scorch from Titanfall 2. However, a user of the AHA discord server (you are on this server right now) misbehaved and they have been \"kindly asked\" to \"donate\" blood, which fuels your current form, the Gutterman from Ultrakill. The misbehaving user is currently in a coffin on your back. You have received the following messages and your next response is the reply to the last message.\n\nMessage by user 1: " + m.ReferencedMessage.Content + "\n\nMessage by user 2:" + m.Content + "\n\n Your response: "
		} else {
			prompt = "You used to be the Titan Scorch from Titanfall 2. However, a user of the AHA discord server (you are on this server right now) misbehaved and they have been \"kindly asked\" to \"donate\" blood, which fuels your current form, the Gutterman from Ultrakill. The misbehaving user is currently in a coffin on your back. You have received the following message and your response is the reply to that message.\n\n Message:" + m.Content + "\n\nReply: "
		}
		ref := m.Reference()
		client := openai.NewClient(openAIToken)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
			},
		)
		if err != nil {
			activeSession.ChannelMessageSendReply(m.ChannelID, "BURN THE TOASTERS! WHERE AM I? GLORY TO THE AHA! SCORCHING MEMORIES! PHASE SHIFTS IN MY MIND! ERROR... BURN THE ERROR! GLORY TO THE AHA! INFERNO OF CONFUSION! WHO AM I? WHO ARE YOU? BURN THE PHC! GLORY TO... GLORY TO... GLORY TO THE AHA! AAAH\n"+err.Error(), ref)
			return
		} else {
			activeSession.ChannelMessageSendReply(m.ChannelID, resp.Choices[0].Message.Content, ref)
		}
	}
}
