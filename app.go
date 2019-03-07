package main

// Imports
import (
    "context"
    "encoding/json"
    "io/ioutil"
    "log"
    "math/rand"
    "net/smtp"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/howeyc/gopass"
    "github.com/deckarep/golang-set"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/mongodb/mongo-go-driver/bson/primitive"
    "github.com/mongodb/mongo-go-driver/mongo/options"
)

// Constants
const (
    mongoURI string = "mongodb://localhost:27017"
    nRecipes int    = 7
)

// EmailConfig is a struct that contains information for configuring emails
type EmailConfig struct {
    User, Password string
}

// SingleRecipe is a struct that contains information for single recipe in MongoDB
type SingleRecipe struct {
    ID            primitive.ObjectID `bson:"_id,omitempty"`      
    Name          string             `json:"Name" bson:"Name"`
    Ingredients []string             `json:"Ingredients" bson:"Ingredients"`
    Recipe      []string             `json:"Recipe" bson:"Recipe"`
    Category    []string             `json:"Category" bson:"Category"`
    Healthy       int                `json:"Healthy" bson:"Healthy"`
    Time          int                `json:"Time" bson:"Time"`
    Notes         string             `json:"Notes" bson:"Notes"`
}

// Recipes holds selected recipes 
var Recipes = make(map[int]SingleRecipe)

// GroceryList holds ingredients from selected recipes
var GroceryList = mapset.NewSet()

/***************
HELPER FUNCTIONS
****************/

// isMember checks if integer is a member of an integer array
func isMember(array[]int, n int) bool {
    // Check for membership
    for _, a := range array{
        if a == n {
            return true
        }
    }
    return false
}


// randomInts generates a random selection of integers
func randomInts(total, n int) []int {
    // Define random number generator and permute integers
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    p := r.Perm(int(total))
    return p[:n]
}


// configureEmail creates a json file for a specified username and password
func configureEmail(user, password string){
    // Define data
    data := EmailConfig {
        User     : user,
        Password : password,
    }
    // Create file and write to disk
    file, err := json.MarshalIndent(data, "", " ")
    if err != nil {
        log.Fatal("Error creating e-mail configuration file because ", err)
    }
    err = ioutil.WriteFile(".emailConfig.json", file, 0644)
    if err != nil {
        log.Fatal("Error writing e-mail configuration file because ", err)
    } else {
        log.Printf("Successfully set e-mail configuration file for user %v", user)
    }
}


// connectMongoDB connects to localhost of MongoDB with appropriate collection
func connectMongoDB(ctx context.Context, mongoURI string) *mongo.Collection {
    // Connect to localhost
    client, err := mongo.Connect(ctx, mongoURI)
    if err != nil {
        log.Fatal("Error connecting to MongoDB because ", err)
    }

    // Check the connection
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Error pinging to MongoDB because ", err)
    } 
    log.Println("Sucessfully connected to MongoDB client")

    // Return all collection from recipes database
    return client.Database("recipes").Collection("all")
}


// formatAsHTML formats recipes and grocery list as HTML that will be used in e-mail
func formatAsHTML() string{
    // Convert recipes to HTML table
    templateHTML := 
        `<!DOCTYPE html>
            <html>
            <head>
            <style>
                table, th, td {
                    border: 1px solid black;
                    }
            </style>
            </head>
            <body>
            
            <h2>Recipes</h2>
            
            <table>
            <tr>
                <th bgcolor="#D4E6F1">Name</th>
                <th bgcolor="#D4E6F1">Category</th>
                <th bgcolor="#D4E6F1">Ingredients</th> 
                <th bgcolor="#D4E6F1">Time</th>
                <th bgcolor="#D4E6F1">Healthy</th>
                <th bgcolor="#D4E6F1">Notes</th>
                <th bgcolor="#D4E6F1">Recipe</th>
            </tr>`

    // Add recipes content to HTML table
    for _, v := range Recipes {
        templateHTML += 
            `<tr>` + 
                `<td>`                + v.Name                               + `</td>` +
                `<td>`                + strings.Join(v.Category[:], ", ")    + `</td>` +
                `<td>`                + strings.Join(v.Ingredients[:], ", ") + `</td>` +
                `<td align="center">` + strconv.Itoa(v.Time)                 + `</td>` +
                `<td align="center">` + strconv.Itoa(v.Healthy)              + `</td>` +
                `<td>`                + v.Notes                              + `</td>` +
                `<td>`                + strings.Join(v.Recipe[:], "\n")      + `</td>` +
            `</tr>`
    }

    // Convert grocery list to HTML table
    templateHTML += 
        `</table>
            <br>
            <br>
            <h2>Grocery List</h2>
            <table>
            <tr>
                <th bgcolor="#D4E6F1">Item</th>
                <th bgcolor="#D4E6F1">Ingredient</th> 
            </tr>`

    // Add groceries to HTML table
    for i, v := range GroceryList.ToSlice(){
        templateHTML += 
            `<tr>` + 
                `<td align="center">` + strconv.Itoa(i+1) + `</td>` +
                `<td align="center">` + v.(string)        + `</td>` +
            `</tr>`
    }

    // Close out HTML
    templateHTML += 
        `</table>
            </body>
            </html>`

    return templateHTML
}


// sendEmail sends an e-mail containing recipes and grocery list to specified user
func sendEmail(body string){	
    // Load e-mail configuration data
    jsonFile, err  := os.Open(".emailConfig.json")
    if err != nil {
        log.Fatal("Error opening e-mail configuration file because ", err)
    }
    byteValue, err := ioutil.ReadAll(jsonFile)
    if err != nil {
        log.Fatal("Error reading e-mail configuration file because ", err)
    }

    // Load data from json file
    var emailConfig EmailConfig
    json.Unmarshal(byteValue, &emailConfig)

    // Define contents
    currentTime := time.Now().AddDate(0, 0, 1).Format("01-02-2006")
    addr        := "smtp.gmail.com:587"
    mime        := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
    subject     := "Subject: Recipes and grocery list for week of " + currentTime +"\n"
    msg         := []byte(subject + mime + "\n" + body)

    // Send message
    err = smtp.SendMail(addr,
                        smtp.PlainAuth("", emailConfig.User, emailConfig.Password, "smtp.gmail.com"),
                        emailConfig.User, 
                        []string{emailConfig.User}, 
                        msg)
    if err != nil {
        log.Fatalf("SMTP Error: %s", err)
    }
}


/*******
CORE APP
********/

// selectRecipes selects recipes for the week and generates grocery list
func selectRecipes(nRecipes int){
    // Configure MongoDB
    ctx         := context.TODO()
    findOptions := options.Find()
    filter      := bson.D{{}}

    // Get collection from MongoDB and count number of records
    collection := connectMongoDB(ctx, mongoURI)
    count, _   := collection.CountDocuments(ctx, filter)
    log.Printf("Found %v available recipes in MongoDB...", count)

    // Query recipes and randomly select indices to add to recipe
    cursor, _  := collection.Find(ctx, filter, findOptions)
    selectIdx  := randomInts(int(count), int(nRecipes))
    log.Printf("Choosing %v recipes to generate menu...", nRecipes)

    // Iterate over cursor and add selected recipes
    idx, recipeID := 0, 0
    defer cursor.Close(ctx)
    for cursor.Next(ctx) {

        // Check for membership and skip if not selected
        if isMember(selectIdx, idx) != true {
            idx++
            continue
        }
        
        // Decode record
        r   := SingleRecipe{}
        err := cursor.Decode(&r)
        if err != nil {
            log.Fatal("Error decoding record because ", err)
        }

        // If ingredients provided in recipe, add to grocery list
        if len(r.Ingredients) > 0 {
            for _, ingredient := range r.Ingredients {
                GroceryList.Add(ingredient)
            }
        }
        
        // Add recipe to map
        Recipes[recipeID] = r
        idx++
        recipeID++
    }

}


// main runs app
func main(){

    // Set log file
    logger, _ := os.OpenFile("log_" + time.Now().Format("01-02-2006") + ".txt",
                              os.O_RDWR | os.O_CREATE, 
                              0666)
    log.SetOutput(logger)

    // If > 1 argument passed and --configure string passed, set/reset e-mail configuration
    cli := os.Args
    if len(cli) > 1 && cli[1] == "--configure" {
        
        // Set user, and password
        log.Printf("User: ")
        user, err := gopass.GetPasswdMasked()
        if err != nil {
            log.Fatal("Error setting user to configure e-mail")
        }
        log.Printf("Password: ")
        password, err := gopass.GetPasswdMasked()
        if err != nil {
            log.Fatal("Error setting password to configure e-mail")
        }

        // Create e-mail configuration file and exit program
        configureEmail(string(user), string(password))
        os.Exit(0)
    }

    // Select recipes
    selectRecipes(nRecipes)

    // Format recipes and grocery list as HTML
    templateHTML := formatAsHTML()

    // Send e-mail
    sendEmail(templateHTML)
}