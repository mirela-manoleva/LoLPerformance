# LoLPerformance - What is it?
A tool which extracts data from your last LoL game and fills it in a performance Excel table.
It divides all of your queues in different sheets and adds some coloring whether a stat is good or bad.
Some stats (like kills or assists) are subjective so we don't touch those.

# Why should you use it?
The primary reason is to have an easy way to check your progress. I've found that sites like [OP.GG](https://www.op.gg/?hl=en_US) or [League of graphs](https://www.leagueofgraphs.com/) don't summorize the data in a manner that I like. It looks good and it gives you an overview of the last few games but you cannot see if you've improved over the season.
Viewing your progress in an excel spreadsheet sure feels weird (and ugly) compared to those sites but it serves a better purpose if your main goal is to improve. The main idea is taken from [Simba ADC](https://www.youtube.com/watch?v=BnhBC9efvrU).

# Instal

1. Checkout the repo:
```
git clone https://github.com/mirela-manoleva/LoLPerformance.git
```
2. Build:
```
cd [GIT REPO FOLDER]
go build
```

# Use

1. Add your user name, tag and region in the `config/user_config.json` file.
```
{
  "region": "[EUW/EUNE]",
  "summonerName": "[USERNAME]",
  "summonerTag": "[TAG]",
  "excelFile": "Improvement.xlsx"
}
```
Note that the app works only for EUW and EUNE. If you need it for any other region create an issue and we will add it.

2. Add your own API key add it to `config/api.key`. If you share your API key with someone make sure you're not using the app at the same time as you might break some Riot API limits.
```
[ENTER YOUR API KEY HERE]
```

3. Run the `.exe` file after each game in order to record it to your spreadsheet. The first time you run the program it will create the `.xlsx` file.

4. Go (pun intended) smurf on the enemies!